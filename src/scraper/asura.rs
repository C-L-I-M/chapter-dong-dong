use log::warn;
use regex::Regex;
use reqwest;
use scraper::{Html, Selector};
use simplelog::info;
use std::iter::zip;

use super::updated_chapter::UpdatedChapter;

static LAST_UPDATED_URL: &str = "https://asuracomic.net/series";
static URL_SERIES_REGEX: &str = r##"series/(?<series_uri>[-a-zA-Z0-9]+)"##;
// class="block text-[13.3px] font-bold">
// block text\-\[13.3px\] font\-bold">(?<series_name>.+?)<
static CHAPTER_REGEX: &str = r##"block text\-\[13.3px\] font\-bold">(?<series_name>.+?)<.*?(Chapter).*?(?<latest_chapter_number>[0-9]+)"##;
static GRID_LAST_UPDATED_SELECTOR: &str =
    r#"div[class="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-5 gap-3 p-4"]"#;

async fn get_html(url: &str) -> String {
    reqwest::get(url).await.unwrap().text().await.unwrap()
}

async fn get_matches_from_html(html: &str) -> Vec<UpdatedChapter> {
    std::fs::write("resources/asura_series.html", html);
    let mut last_updated_series = vec![];
    let links_regex = Regex::new(URL_SERIES_REGEX).unwrap();
    let name_and_chapter_regex = Regex::new(CHAPTER_REGEX).unwrap();
    for (series_uri_match, name_and_latest_chapter_match) in zip(
        links_regex.captures_iter(html),
        name_and_chapter_regex.captures_iter(html),
    ) {
        let series_uri = format!(
            "{}/{}",
            LAST_UPDATED_URL,
            series_uri_match.name("series_uri").unwrap().as_str(),
        );
        let series_name = name_and_latest_chapter_match
            .name("series_name")
            .unwrap()
            .as_str();
        let chapter_number = name_and_latest_chapter_match
            .name("latest_chapter_number")
            .unwrap()
            .as_str();
        let chapter = UpdatedChapter {
            series_name: String::from(series_name),
            series_uri: String::from(series_uri),
            latest_chapter: chapter_number.parse::<u32>().unwrap(),
        };
        last_updated_series.push(chapter)
    }
    last_updated_series
}

async fn parse_html(html: &str) -> String {
    let selector = Selector::parse(GRID_LAST_UPDATED_SELECTOR).unwrap();
    Html::parse_document(html)
        .select(&selector)
        .next()
        .unwrap()
        .inner_html()
}

pub async fn get_last_updated_series() -> Vec<UpdatedChapter> {
    let raw_html = get_html(LAST_UPDATED_URL).await;
    let last_updated_series = get_matches_from_html(&parse_html(&raw_html).await).await;

    warn!("found {} matches", last_updated_series.len());
    last_updated_series
}

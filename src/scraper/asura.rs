use regex::Regex;
use reqwest;
use scraper::{Html, Selector};

static LAST_UPDATED_URL: &str = "https://asuracomic.net/series";
static REGEX: &str = r##"series/(?<series_name>[-a-zA-Z0-9]+)"##;
static GRID_LAST_UPDATED_SELECTOR: &str =
    r#"div[class="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-5 gap-3 p-4"]"#;

async fn get_html(url: &str) -> String {
    reqwest::get(url).await.unwrap().text().await.unwrap()
}

async fn get_matches_from_html(html: &str) -> Vec<String> {
    let mut last_updated_series: Vec<String> = vec![];
    let re = Regex::new(REGEX).unwrap();
    for m in re.captures_iter(html) {
        if let Some(name) = m.name("series_name") {
            last_updated_series.push(String::from(name.as_str()));
        }
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

pub async fn get_last_updated_series() -> Vec<String> {
    let raw_html = get_html(LAST_UPDATED_URL).await;
    let last_updated_series: Vec<String> =
        get_matches_from_html(&parse_html(&raw_html).await).await;
    for m in &last_updated_series {
        println!("{}", m);
    }
    println!("found {} matches", last_updated_series.len());
    last_updated_series
}

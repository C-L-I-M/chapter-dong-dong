use crate::discord;
use crate::scraper::updated_chapter::UpdatedChapter;
use serde_derive::{Deserialize, Serialize};
use simplelog::warn;
use std::collections::HashMap;
use std::fs;
use std::fs::File;

#[derive(Deserialize, Serialize)]
struct Serie {
    pub name: String,
    pub latest_chapter: u32,
}

pub async fn update_latest_series(store_path: &str, updated_series: Vec<UpdatedChapter>) {
    let data_stored =
        std::fs::read_to_string(store_path).expect(&format!("No such file: {}", store_path));
    let mut json_stored: HashMap<String, Serie> = serde_json::from_str(&data_stored).unwrap();
    for series in updated_series {
        if let Some(entry) = json_stored.get_mut(&series.series_uri) {
            if entry.latest_chapter == series.latest_chapter {
                continue;
            }
            warn!(
                "Updated {} to chapter {}",
                series.series_name, series.latest_chapter
            );
            entry.latest_chapter = series.latest_chapter;

            let chapter_uri = format!("{}/chapter/{}", series.series_uri, series.latest_chapter);
            let message = format!(
                "Chapter {} of `{}` is available!\n{}",
                series.latest_chapter, series.series_name, chapter_uri
            );
            discord::send_message(&message).await
        } else {
            warn!(
                "Initialized {} to chapter {}",
                series.series_name, series.latest_chapter
            );

            let to_insert = Serie {
                name: String::from(&series.series_name),
                latest_chapter: series.latest_chapter,
            };

            json_stored.insert(String::from(&series.series_uri), to_insert);
            let message = format!(
                "New series available: `{}`, {} chapters\n{}",
                &series.series_name, series.latest_chapter, &series.series_uri
            );
            discord::send_message(&message).await
        }
    }

    fs::write(store_path, serde_json::to_string(&json_stored).unwrap());
}

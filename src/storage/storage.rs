use crate::scraper::updated_chapter::UpdatedChapter;
use serde_derive::{Deserialize, Serialize};
use simplelog::info;
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
            info!(
                "Updated {} to chapter {}",
                series.series_name, series.latest_chapter
            );
            entry.latest_chapter = series.latest_chapter
        } else {
            info!(
                "Initialized {} to chapter {}",
                series.series_name, series.latest_chapter
            );

            let to_insert = Serie {
                name: series.series_name,
                latest_chapter: series.latest_chapter,
            };

            json_stored.insert(String::from(series.series_uri), to_insert);
        }
    }

    info!("Saving data to store");
    fs::write(store_path, serde_json::to_string(&json_stored).unwrap());
}

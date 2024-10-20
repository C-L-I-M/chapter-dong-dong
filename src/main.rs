use simplelog::*;

use chapter_dong_dong::cli;
use chapter_dong_dong::discord;
use chapter_dong_dong::scraper;
use chapter_dong_dong::storage;
use std::{thread, time};
use tokio::runtime::Handle;

async fn test_entrypoint() -> () {
    scraper::asura::get_last_updated_series().await;
}

async fn poll_asura(store_path: String) -> () {
    warn!("!!!!");
    loop {
        warn!("getting update from asura");
        storage::storage::update_latest_series(
            &store_path,
            scraper::asura::get_last_updated_series().await,
        )
        .await;
        warn!("sleeping for 1 minute");
        thread::sleep(time::Duration::from_secs(60));
    }
}

#[tokio::main]
async fn main() {
    CombinedLogger::init(vec![TermLogger::new(
        LevelFilter::Warn,
        Config::default(),
        TerminalMode::Mixed,
        ColorChoice::Auto,
    )])
    .unwrap();

    let cli = cli::build_cli();
    let matches = cli.get_matches();

    let store_path = if let Some(value) = matches.get_one::<String>("store") {
        String::from(value)
    } else {
        String::from("s3_storage.json")
    };

    if matches.get_flag("debug_entrypoint") {
        test_entrypoint().await;
    }

    let mut client = match discord::build_client().await {
        Ok(client) => client,
        Err(err) => panic!("{err:?}"),
    };

    if let Some(value) = matches.get_one::<String>("msg") {
        warn!("Send test message: {}", value);
        discord::send_message(value).await
    }

    let h = Handle::current();
    h.spawn(async move {
        warn!("Starting asura polling");
        poll_asura(String::from(store_path)).await;
    });
    h.spawn(async move {
        warn!("Starting bot");
        client.start().await;
    });
    warn!("Started everying, wait for the dong-dong!");
    loop {}
}

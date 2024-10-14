use serenity::http::Http;
use simplelog::*;

use chapter_dong_dong::cli;
use chapter_dong_dong::discord;
use chapter_dong_dong::scraper;

async fn test_entrypoint() -> () {
    scraper::asura::get_last_updated_series().await;
}

#[tokio::main]
async fn main() {
    CombinedLogger::init(
        vec![
            TermLogger::new(LevelFilter::Warn, Config::default(), TerminalMode::Mixed, ColorChoice::Auto),
        ]
    ).unwrap();

    let cli = cli::build_cli();
    let matches = cli.get_matches();


    if matches.get_flag("debug_entrypoint") {
        test_entrypoint().await;
    }

    let mut client = match discord::build_client().await {
        Ok(client) => client,
        Err(err) => panic!("{err:?}")
    };

    if let Some(value) = matches.get_one::<String>("msg") {
        info!("Send test message: {}", value);
        discord::send_message(value).await
    }

    let future = client.start();

    future.await.expect("TODO: panic message");
}

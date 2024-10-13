use chapter_dong_dong::cli;
use chapter_dong_dong::discord;
use chapter_dong_dong::scraper;

async fn start_bot() {
    let client = discord::build_client().await;
    let client = match client {
        Ok(mut unwrapped_client) => unwrapped_client.start().await,
        Err(_) => panic!("ok!"),
    };
    match client {
        Ok(_current) => {
            println!("Starting the dong-dong!")
        }
        Err(err) => panic!("{err:?}"),
    }
}

async fn test_entrypoint() -> () {
    scraper::asura::get_last_updated_series().await;
}

#[tokio::main]
async fn main() {
    let cli = cli::build_cli();
    let matches = cli.get_matches();
    if matches.get_flag("debug_entrypoint") {
        test_entrypoint().await;
    } else {
        start_bot().await;
    }
}

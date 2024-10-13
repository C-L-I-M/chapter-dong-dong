use chapter_dong_dong::cli;
use chapter_dong_dong::discord::build_client;

async fn start_bot() {
    let client = build_client().await;
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
    println!("noop")
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

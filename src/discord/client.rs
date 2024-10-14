use poise::serenity_prelude::{ClientBuilder, GatewayIntents};
use serenity::all::{ChannelId, Http};

use crate::discord::command;
use crate::discord::types::Data;

fn discord_token() -> String {
    std::env::var("DISCORD_TOKEN").expect("missing DISCORD_TOKEN")
}

pub async fn send_message(msg: &str) {
    let http = &Http::new(&discord_token());
    ChannelId::new(1293537672287227924u64).say(http, msg).await;
}

pub async fn build_client() -> Result<poise::serenity_prelude::Client, serenity::prelude::SerenityError> {
    let token = discord_token();
    let intents = GatewayIntents::non_privileged();

    let framework = poise::Framework::builder()
        .options(poise::FrameworkOptions {
            commands: vec![command::age()],
            ..Default::default()
        })
        .setup(|ctx, _ready, framework| {
            Box::pin(async move {
                poise::builtins::register_globally(ctx, &framework.options().commands).await?;
                Ok(Data {})
            })
        })
        .build();

    ClientBuilder::new(token, intents)
        .framework(framework)
        .await
}

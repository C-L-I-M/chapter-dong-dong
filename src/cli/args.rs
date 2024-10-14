use clap::{arg, Command};

pub fn build_cli() -> Command {
    Command::new("CLI")
        .arg(arg!(-d - -debug_entrypoint).action(clap::ArgAction::SetTrue))
        .arg(arg!(-m - -msg).action(clap::ArgAction::Set))
}

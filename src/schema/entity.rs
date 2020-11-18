use super::record::Record;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Entity;

impl Record for Entity {
    fn id(&self) -> String {
        return String::from("hoi");
    }
    fn ident(&self) -> String {
        return String::from("hoi");
    }
    fn set_ident(&mut self, ident: String) {}
}

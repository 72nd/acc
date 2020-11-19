use super::record::{Record, RecordType};

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Project;

impl Record for Project {
    fn id(&self) -> String {
        return String::from("hoi");
    }
    fn ident(&self) -> String {
        return String::from("hoi");
    }
    fn set_ident(&mut self, ident: String) {}
    fn record_type(&self) -> RecordType {
        RecordType::Entity
    }
}

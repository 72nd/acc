use super::common::{Ident, ID};
use super::record::{Record, RecordType};

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Project;

impl Record for Project {
    fn id(&self) -> ID {
        return ID::new();
    }
    fn ident(&self) -> Ident {
        return Ident::from("p-23");
    }
    fn set_ident(&mut self, ident: Ident) {}
    fn record_type(&self) -> RecordType {
        RecordType::Project
    }
}

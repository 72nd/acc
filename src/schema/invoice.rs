use super::common::{ID, Ident};
use super::record::{Record, RecordType};

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Invoice;

impl Record for Invoice {
    fn id(&self) -> ID {
        return ID::new();
    }
    fn ident(&self) -> Ident {
        return Ident::from("i-23")
    }
    fn set_ident(&mut self, ident: Ident) {}
    fn record_type(&self) -> RecordType {
        RecordType::Transaction
    }
}

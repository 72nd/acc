use super::common::ID;
use super::record::{Record, RecordType};

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Transaction;

impl Record for Transaction {
    fn id(&self) -> ID {
        return ID::new();
    }
    fn ident(&self) -> String {
        return String::from("hoi");
    }
    fn set_ident(&mut self, ident: String) {}
    fn record_type(&self) -> RecordType {
        RecordType::Transaction
    }
}

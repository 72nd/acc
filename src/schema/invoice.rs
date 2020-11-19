use super::common::ID;
use super::record::{Record, RecordType};

use serde::{Deserialize, Serialize};

/// The identifier prefix for Invoices.
const INVOICE_IDENT_PREFIX: &str = "t";

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Invoice;

impl Record for Invoice {
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

use super::record::{Record, RecordType};

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExpenseCategory;

impl Record for ExpenseCategory {
    fn id(&self) -> String {
        return String::from("hoi");
    }
    fn ident(&self) -> String {
        return String::from("hoi");
    }
    fn set_ident(&mut self, ident: String) {}
    fn record_type(&self) -> RecordType {
        RecordType::ExpenseCategory
    }
}

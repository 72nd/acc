use super::common::ID;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
/// Expense represents a payment done by the company or a third party to assure the ongoing of
/// the business.
pub struct Expense {
    /// Internal unique identifier of the Expense.
    id: ID,
    /// User chosen and human readable identifier. This is helpful to mark a record and it's
    /// attachments more understandable as only using 
    ident: String,
    /// Describes the Expense in a meaningful way.
    name: String,

}

impl Expense {
}

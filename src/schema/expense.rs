use super::common::ID;
use super::date::Date;
use super::entity::Entity;
use super::money::Money;
use super::record::{Record, RecordType};
use super::relation::Relation;

use std::path::PathBuf;

use serde::{Deserialize, Serialize};

#[derive(Debug, Clone)]
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
    /// States the amount payed for the Expense.
    amount: Money,
    /// Path to an attachment which serves as a business record for this expense. The most common
    /// example is the scan of an invoice or a receipt.
    path: PathBuf,
    /// The day the obligation of this expense emerged.
    date_of_accrual: Date,
    /// States whether the Expense has to be forwarded to any third party.
    billable: bool,
    /// Refers to the customer which has to pay for the expense (if this is the case).
    obliged_customer: Relation<Entity>,
    /// States whether any third party (most common: a employee) has advanced this expense for the
    /// company and needs a pay back.
    advanced_by_third_party: bool,
    /// Refers to the third party which advanced the payment.
    advanced_third_party: Relation<Entity>,
    /// The date of the settlement of the expense (the company has not to take further actions).
    date_of_settlement: Date,
}

impl Record for Expense {
    fn id(&self) -> String {
        return String::from("hoi");
    }
    fn ident(&self) -> String {
        return String::from("hoi");
    }
    fn set_ident(&mut self, ident: String) {}
    fn record_type(&self) -> RecordType {
        RecordType::Expense
    }
}

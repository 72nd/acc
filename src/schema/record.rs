use super::common::{ID, Ident};
use super::expense::Expense;

use std::fmt;

/// A collection of multiple Records. This structure bundles all common and repeating task which
/// are not specific to a certain type of Element collection.
pub struct Records<R: Record>(Vec<R>);

/// Records are the fundamental building block of Acc. They represent one data entry in the system.
/// Examples: Expense, Customer. As most of the operations on the different collections of records
/// are the same, these are combined using this trait.
pub trait Record {
    /// Return the ID of a Record.
    fn id(&self) -> ID;
    /// Return the Identifier of a Record.
    fn ident(&self) -> Ident;
    /// Set the Identifier of a Record. As Identifiers have to be unique, it's important to set
    /// this value always trough the Records type.
    fn set_ident(&mut self, ident: Ident);
    /// Returns the type of the record.
    fn record_type(&self) -> RecordType;
}

/// The enumeration describes the different types of records existing. Each type is described in
/// it's structure and implements the Record interface. This enumeration is mainly used to provide
/// more helpful (debug) messages for the user.
#[derive(Debug, Clone)]
pub enum RecordType {
    Customer,
    Employee,
    Expense,
    ExpenseCategory,
    Invoice,
    Misc,
    Project,
    Transaction,
}

impl fmt::Display for RecordType {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        f.write_str(match self {
            Customer => "Customer",
            Employee => "Employee",
            Expense => "Expense",
            ExpenseCategory => "Expense Category",
            Invoice => "Invoice",
            Misc => "Miscellaneous Record",
            Project => "Project",
            Transaction => "Transaction",
        })
    }
}

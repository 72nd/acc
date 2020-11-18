use super::Expense;

use std::fmt;

/// A collection of multiple Records. This structure bundles all common and repeating task which
/// are not specific to a certain type of Element collection.
pub struct Records<R: Record>(Vec<R>);

/// Records are the fundamental building block of Acc. They represent one data entry in the system.
/// Examples: Expense, Customer. As most of the operations on the different collections of records
/// are the same, these are combined using this trait.
pub trait Record {
    /// Return the ID of a Record.
    fn id(&self) -> String;
    /// Return the Identifier of a Record.
    fn ident(&self) -> String;
    /// Set the Identifier of a Record. As Identifiers have to be unique, it's important to set
    /// this value always trough the Records type.
    fn set_ident(&mut self, ident: String);
}

/// The enumeration describes the different types of records existing. Each type is described in
/// it's structure and implements the Record interface. This enumeration is mainly used to provide
/// more helpful (debug) messages for the user.
#[derive(Debug, Clone)]
pub enum RecordType {
    Entity,
    Expense,
    Invoice,
    Misc,
    Project,
    Transaction,
}

impl fmt::Display for RecordType {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        f.write_str(match self {
            Entity => "Entity",
            Expense => "Expense",
            Invoice => "Invoice",
            Misc => "Miscellaneous Record",
            Project => "Project",
            Transaction => "Transaction",
        })
    }
}

/// A Relation is used to describe a relation between two Records. Currently relations only
/// describes One-To-One relations.
#[derive(Debug, Clone)]
pub struct Relation<T: Record>(T);

impl<T: Record> Relation<T> {
    pub fn element(&self) -> T {
        return self.0;
    }
}

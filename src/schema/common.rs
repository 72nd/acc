use serde::{Deserialize, Serialize};
use uuid::Uuid;

/// A collection of multiple Records. This structure bundles all common and repeating task which
/// are not specific to a certain type of Element collection.
pub struct Records<R: Record>(Vec<R>);

/// Records are the fundamental building block of Acc. They represent one data entry in the system.
/// Examples: Expense, Customer. As most of the operations on the different collections of records
/// are the same, these are combined using this trait.
pub trait Record {
    /// Return the ID of a Record.
    fn id() -> String;
    /// Return the Identifier of a Record.
    fn ident() -> String;
    /// Set the Identifier of a Record. As Identifiers have to be unique, it's important to set
    /// this value always trough the Records type.
    fn set_ident(&mut self, ident: String);
}

/// System wide UUID to identify all records. Using UUID Version 4.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ID(Uuid);

impl ID {
    /// Returns a new instance of the ID with a freshly generated UUID.
    fn new() -> Self {
        Self(Uuid::new_v4())
    }
}

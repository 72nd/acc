use super::common::{Ident, ID};
use super::record::{Record, RecordType};

use serde::{Deserialize, Serialize};

/// The type of en entity. Can be a customer or a employee.
pub trait EntityType {
    /// Returns the Record Type of the specific entity (Customer or Employee).
    fn record_type(&self) -> RecordType;
}

/// Any third-party which is a customer of the company.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Customer {}
impl EntityType for Customer {
    fn record_type(&self) -> RecordType {
        RecordType::Customer
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
/// Somebody which works for the company.
pub struct Employee {}

impl EntityType for Employee {
    fn record_type(&self) -> RecordType {
        RecordType::Employee
    }
}

/// A entity is a person or a company which has some relation with the company. Entities are used
/// to depict employees, customers and so on. This type was called "Party" in the go version of
/// Acc.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Entity<T>
where
    T: EntityType,
{
    /// Unique internal identifier of the Expense.
    id: ID,
    /// User-chosen and human readable identifier. This is helpful to mark a record and it's
    /// attachments more understandable as only using some long UUID.
    ident: Ident,
    /// The name of the entity. For natural persons this should contain the full name (first and
    /// last name).
    name: String,
    /// The street name with number and any additional delivery instructions.
    street: String,
    /// The postal aka ZIP-code of the entities address.
    postal_code: String,
    /// The place of the entities address.
    place: String,
    /// Gives information about the type of entity.
    entity_type: T,
}

impl<T> Record for Entity<T>
where
    T: EntityType,
{
    fn id(&self) -> ID {
        return self.id;
    }
    fn ident(&self) -> Ident {
        return self.ident.clone();
    }
    fn set_ident(&mut self, ident: Ident) {}
    fn record_type(&self) -> RecordType {
        self.entity_type.record_type()
    }
}

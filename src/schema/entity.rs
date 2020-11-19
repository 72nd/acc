use super::common::ID;
use super::record::{Record, RecordType};

use serde::{Deserialize, Serialize};

/// Describes the type of entity.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum EntityType {
    /// A person or company who buys stuff from the company.
    Customer,
    /// Someone who works for the company itself.
    Employee,
}

/// A entity is a person or a company which has some relation with the company. Entities are used
/// to depict employees, customers and so on. This type was called "Party" in the go version of
/// Acc.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Entity {
    /// Unique internal identifier of the Expense.
    id: ID,
    /// User-chosen and human readable identifier. This is helpful to mark a record and it's
    /// attachments more understandable as only using some long UUID.
    ident: String,
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
    entity_type: EntityType,
}

impl Record for Entity {
    fn id(&self) -> ID {
        return self.id;
    }
    fn ident(&self) -> String {
        return self.ident.clone();
    }
    fn set_ident(&mut self, ident: String) {}
    fn record_type(&self) -> RecordType {
        RecordType::Entity
    }
}

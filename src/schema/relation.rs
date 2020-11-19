use super::common::ID;
use super::record::Record;

use std::fmt;

use serde::de::{self, Deserialize, Deserializer};
use serde::{Serialize, Serializer};

/// A Relation is used to describe a relation between two Records. Currently relations only
/// describes One-To-One relations. The relation is expressed in the UUID of the record.
#[derive(Debug, Clone, PartialEq)]
pub struct Relation<T: Record> {
    /// The Id which links to the referenced record.
    id: ID,
    /// Content of the referenced object.
    value: Option<T>,
}

impl<T: Record> Relation<T> {
    /// Takes a ID and returns a new instance of the Relation structure.
    pub fn new(id: ID) -> Self {
        Self {
            id: id,
            value: None,
        }
    }
    /*
    pub fn element(&self) -> T {
        return self.0;
    }
    */
}

impl<T: Record> fmt::Display for Relation<T> {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match &self.value {
            Some(x) => write!(f, "relation to {} {}", self.id, x.record_type()),
            None => write!(f, "relation to record {}", self.id),
        }
    }
}

impl<T: Record> Serialize for Relation<T> {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        serializer.serialize_str(&self.id.to_string())
    }
}

impl<'de, T: Record> Deserialize<'de> for Relation<T> {
    fn deserialize<D>(deserializer: D) -> Result<Relation<T>, D::Error>
    where
        D: Deserializer<'de>,
    {
        let value: String = Deserialize::deserialize(deserializer)?;
        match ID::from(&value) {
            Ok(x) => Ok(Relation::new(x)),
            Err(_) => Err(de::Error::custom(format!(
                "given input \"{}\" couldn't be parsed as a UUID for a relation",
                value
            ))),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::Relation;
    use super::ID;
    use crate::schema::{Customer, Entity};

    #[test]
    fn test_relation_serialize() {
        let relation = Relation::<Entity<Customer>> {
            id: ID::new(),
            value: None,
        };
        let expected = format!("---\n{}", relation.id);
        let serialized = serde_yaml::to_string(&relation).unwrap();
        assert_eq!(expected, serialized);
    }

    #[test]
    fn test_relation_deserialize() {
        let expected = Relation::<Entity<Customer>> {
            id: ID::new(),
            value: None,
        };
        let input = format!("---\n{}", expected.id);
        assert_eq!(
            expected.id,
            serde_yaml::from_str::<Relation<Entity<Customer>>>(&input)
                .unwrap()
                .id
        );
    }
}

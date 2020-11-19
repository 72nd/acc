use std::fmt;

use serde::{Deserialize, Serialize};
use uuid::{Error as UuidError, Uuid};

/// System wide UUID to identify all records. Using UUID Version 4.
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub struct ID(Uuid);

impl ID {
    /// Returns a new instance of the ID with a freshly generated UUID.
    pub fn new() -> Self {
        Self(Uuid::new_v4())
    }

    /// Tries to parse the UUID from a given string and returns a ID object.
    pub fn from(input: &str) -> Result<Self, UuidError> {
        match Uuid::parse_str(input) {
            Ok(x) => Ok(Self(x)),
            Err(e) => Err(e),
        }
    }
}

impl fmt::Display for ID {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

use serde::{Deserialize, Serialize};
use uuid::Uuid;


/// System wide UUID to identify all records. Using UUID Version 4.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ID(Uuid);

impl ID {
    /// Returns a new instance of the ID with a freshly generated UUID.
    fn new() -> Self {
        Self(Uuid::new_v4())
    }
}

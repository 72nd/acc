use super::date::Date;
use super::record::{RecordType, Typer};

use std::fmt;

use serde::{Deserialize, Serialize};
use uuid::{Error as UuidError, Uuid};

/// System wide UUID to identify all records. Using UUID Version 4.
#[derive(Debug, Clone, Copy, Serialize, Deserialize, PartialEq)]
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

/// Short but unique identifier for a record. The purpose of this identifier is to provide the
/// user a more memorable id than the internally used UUID of the ID object. For some users this
/// will easy the integration of Acc within any existing work flows as this field has no
/// requirements to it's content besides it has to be unique.
///
/// The structure provides some extra functionality to ease the handling with the ident system of
/// the Solutionsbüro.
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub struct Ident<T>
where
    T: Typer,
{
    /// The Typer of the Record which is described.
    record: T,
    /// The actual Identifier.
    data: String,
}

impl<T> Ident<T>
where
    T: Typer,
{
    /// Returns a new Identifier based on a freely-chosen string.
    pub fn from<S: Into<String>, U>(record: T, ident: S) -> Self {
        Self {
            record: record,
            data: ident.into(),
        }
    }

    /// Returns a new Identifier containing the appropriate prefix followed by the given number. As
    pub fn from_n(record: T, number: u64) -> Self {
        Self {
            record: record,
            data: format!("{}-{}", Ident::prefix(), number),
        }
    }

    /// Returns a new Identifier containing the appropriate prefix followed by the given year (two
    /// digit representation) and number.
    pub fn from_year_n(record: T, date: Date, number: u64) -> Self {
        Self {
            record: record,
            data: format!("{}-{}-{}", Ident::prefix(), date.year_2d(), number),
        }
    }

    /// Returns the Solutionsbüro default prefix for any given record type. Each type has it's
    /// individual prefix to clearly distinguish between the different types of records. As this is
    /// the default implementation it returns `?`.
    fn prefix<'a>() -> &'a str {
        "?"
    }

    /*
    fn prefix_by_type<'a>(rtype: RecordType) -> &'a str {
        match rtype {
            RecordType::Customer => "c",
            RecordType::Employee => "y",
            RecordType::Expense => "e",
            RecordType::ExpenseCategory => "ec",
            RecordType::Invoice => "i",
            RecordType::Misc => "m",
            RecordType::Project => "p",
            RecordType::Transaction => "t",
        }
    }
    */
}

impl fmt::Display for Ident {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        f.write_str(&self.0)
    }
}

impl Default for Ident {
    fn default() -> Self {
        Ident::from("e-23")
    }
}

#[cfg(test)]
mod tests {
    use super::Ident;
    use crate::schema::Date;
    use crate::schema::RecordType;

    #[test]
    fn test_from_n() {
        assert_eq!("e-23", Ident::from_n(RecordType::Expense, 23).to_string());
    }

    #[test]
    fn test_from_year_n() {
        let date = Date::from_ymd(2019, 10, 2);
        assert_eq!(
            "e-19-23",
            Ident::from_year_n(RecordType::Expense, date, 23).to_string()
        );
    }
}

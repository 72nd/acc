use std::fmt;

use chrono::NaiveDate;
use serde::de::{self, Deserialize, Deserializer};
use serde::{Serialize, Serializer};

/// Date format used to save values to YAML files. This complies with ISO 8601.
const YAML_DATE_FORMAT: &str = "%Y-%m-%d";

/// Handles dates as a calendar day. Encapsulates the `chrono::NaiveDate` structure and adds custom
/// formatting and Serde (de)serialization. For saving the date is formatted according to ISO 8601
/// (YYYY-MM-DD).
///
/// It's possible to introduce some locale specific behavior down the road. Also the
/// time-zone (which is ignored in the moment) should be addressed in the future.
#[derive(Clone, Debug, PartialEq)]
pub struct Date(NaiveDate);

impl Date {
    pub fn from(date: NaiveDate) -> Self {
        Self(date)
    }
}

impl fmt::Display for Date {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0.format("%a%e. %B %Y"))
    }
}

impl Serialize for Date {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        serializer.serialize_str(&self.0.format(YAML_DATE_FORMAT).to_string())
    }
}

impl<'de> Deserialize<'de> for Date {
    fn deserialize<D>(deserializer: D) -> Result<Date, D::Error>
    where
        D: Deserializer<'de>,
    {
        let value: String = Deserialize::deserialize(deserializer)?;
        match NaiveDate::parse_from_str(&value, YAML_DATE_FORMAT) {
            Ok(x) => Ok(Date::from(x)),
            Err(_) => Err(de::Error::custom(format!(
                "given input value \"{}\" couldn't be parsed as a date in the format {}",
                value, YAML_DATE_FORMAT
            ))),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::Date;

    use chrono::prelude::*;

    #[test]
    fn test_fmt_display() {
        let date = Date::from(NaiveDate::from_ymd(2020, 10, 2));
        assert_eq!("Fri 2. October 2020", format!("{}", date));
    }

    #[test]
    fn test_date_serialize() {
        let date = Date::from(NaiveDate::from_ymd(2020, 10, 2));
        let serialized = serde_yaml::to_string(&date).unwrap();
        assert_eq!("---\n2020-10-02", serialized);
    }

    #[test]
    fn test_date_deserialize() {
        let expected = Date::from(NaiveDate::from_ymd(2020, 10, 2));
        let input = "---\n2020-10-02";
        assert_eq!(expected, serde_yaml::from_str(&input).unwrap());
    }
}

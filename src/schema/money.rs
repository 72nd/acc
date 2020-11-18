use std::fmt;

use rusty_money::{Formatter, Money as RMoney, Params, Position};
use serde::de::{self, Deserialize, Deserializer};
use serde::{Serialize, Serializer};

/// Save representation of a amount of money containing the currency. Money uses the rusty_money
/// crate to handling money amounts.
///
/// It expands the library with the Serde (de-)serialize in our in-house money representation.
/// This format is described in the `Money::yaml_params` method.
#[derive(Eq, PartialEq)]
pub struct Money(RMoney);

impl Money {
    /// Returns a new instance of the Money element. Needs an rusty_money Money instance.
    pub fn new(money: RMoney) -> Self {
        Self(money)
    }

    /// Returns the money Formatter for saving a amount into YAML file. As different currencies
    /// have different formatting, orderings etc. it's important to define one form for the saving.
    ///
    /// - The amount precedences the currency code.
    /// - The minor units are separated by a dot. The rest of the amount is written without any
    /// additional characters.
    /// - The currency is separated by a single whitespace from the amount.
    /// - Currencies are expressed according to ISO 4217.
    fn yaml_params(&self) -> Params {
        let currency = self.0.currency();
        Params {
            // This char has to be removed before saving it.
            digit_separator: '_',
            exponent_separator: '.',
            separator_pattern: vec![3, 3, 3],
            positions: vec![
                Position::Sign,
                Position::Amount,
                Position::Space,
                Position::Code,
            ],
            rounding: Some(currency.exponent),
            symbol: Some(currency.symbol),
            code: Some(currency.iso_alpha_code),
        }
    }
}

impl fmt::Debug for Money {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl fmt::Display for Money {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.0)
    }
}

impl Serialize for Money {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut rsl = Formatter::money(&self.0, self.yaml_params());
        rsl.retain(|c| c != '_');
        serializer.serialize_str(&rsl)
    }
}

impl<'de> Deserialize<'de> for Money {
    fn deserialize<D>(deserializer: D) -> Result<Money, D::Error>
    where
        D: Deserializer<'de>,
    {
        let value: String = Deserialize::deserialize(deserializer)?;

        let parts: Vec<&str> = value.split(' ').collect();
        if parts.len() != 2 {
            return Err(de::Error::custom(format!("given input value \"{}\" couldn't be parsed as a amount of money, contains to many whitespaces", value)));
        }
        match RMoney::from_str(parts[0], parts[1]) {
            Ok(x) => Ok(Money::new(x)),
            Err(e) => Err(de::Error::custom(format!(
                "given input value \"{}\" couldn't be parsed as a amount of money, {}",
                value, e
            ))),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::Money;

    use rusty_money::Money as RMoney;

    #[test]
    fn test_money_serialize() {
        let money = Money::new(RMoney::from_str("2000", "CHF").unwrap());
        let serialized = serde_yaml::to_string(&money).unwrap();
        assert_eq!("---\n2000.00 CHF", serialized);
    }

    #[test]
    fn test_money_deserialize() {
        let expected = Money::new(RMoney::from_str("2000", "CHF").unwrap());
        let input = "---\n2000.00 CHF";
        assert_eq!(expected, serde_yaml::from_str(&input).unwrap())
    }
}

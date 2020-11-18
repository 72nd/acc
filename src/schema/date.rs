use std::fmt;

use chrono::{DateTime, Utc};
use serde::{Serialize, Deserialize};

/// Handles dates as a calendar day. Encapsulates the `chrono::DateTime` structure and hides any
/// time-related function. The usage of the time makes no sense for Acc.
#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Date(DateTime<Utc>);


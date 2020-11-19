use super::common::{Ident, ID};
use super::expense::Expense;

use std::fmt;

/// A collection of multiple Records. This structure bundles all common and repeating task which
/// are not specific to a certain type of Element collection.
pub struct Records<R: Record>(Vec<R>);

/// Records are the fundamental building block of Acc. They represent one data entry in the system.
/// Examples: Expense, Customer. As most of the operations on the different collections of records
/// are the same, these are combined using this trait.
pub trait Record {
    /// Return the ID of a Reord.
    fn id(&self) -> ID;
    /// Return the Identifier of a Record.
    fn ident(&self) -> Ident;
    /// Set the Identifier of a Record. As Identifiers have to be unique, it's important to set
    /// this value always trough the Records type.
    fn set_ident(&mut self, ident: Ident);
    /// Returns the type of the record.
    fn record_type(&self) -> RecordType;
}

/// The enumeration describes the different types of records existing. Each type is described in
/// it's structure and implements the Record interface. This enumeration is mainly used to provide
/// more helpful (debug) messages for the user.
#[derive(Debug, Clone)]
pub enum RecordType {
    Customer,
    Employee,
    Expense,
    ExpenseCategory,
    Invoice,
    Misc,
    Project,
    Transaction,
}

impl fmt::Display for RecordType {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        f.write_str(match self {
            Customer => "Customer",
            Employee => "Employee",
            Expense => "Expense",
            ExpenseCategory => "Expense Category",
            Invoice => "Invoice",
            Misc => "Miscellaneous Record",
            Project => "Project",
            Transaction => "Transaction",
        })
    }
}

/// A collection of Requirements. Normally contains all Checks for one Record type.
pub struct Check<T>
where
    T: Fn() -> bool,
{
    /// The name of the Check/the Record type it runs on to provide additional information to the
    /// user where a Requirement isn't met yet.
    record_name: String,
    /// Contains the Requirements to be tested.
    requirements: Vec<Requirement<T>>,
}

impl<T> Check<T>
where
    T: Fn() -> bool,
{
    /// Returns a new instance of the Check class.
    pub fn new(record_name: String) -> Self {
        Self {
            record_name: record_name,
            requirements: Vec::<Requirement<T>>::new(),
        }
    }

    /// Add a Requirement to the Check.
    pub fn add(&mut self, r: Requirement<T>) -> &mut Self {
        self.requirements.push(r);
        self
    }

    /// Run the check.
    pub fn run(&self) -> CheckResult {
        let mut rsl = CheckResult::new(&self.record_name);
        for r in &self.requirements {
            match &r.check {
                Some(c) => match c() {
                    true => {}
                    false => rsl.add(r),
                },
                None => rsl.log_no_check(),
            };
        }
        rsl
    }
}

/// The results of a Check run. Contains the explanations for the user for each failed requirement
/// check. The Check Results are in a separate to enable more complex feedback to the user in the
/// future.
pub struct CheckResult {
    /// The name of the Check/the Record type it runs on to provide additional information to the
    /// user where a Requirement isn't met yet.
    record_name: String,
    /// Contains all on_fail messages.
    messages: Vec<String>,
}

impl CheckResult {
    /// Returns a new instance of the CheckResult.
    pub fn new<S: Into<String>>(record_name: S) -> Self {
        Self {
            record_name: record_name.into(),
            messages: Vec::<String>::new(),
        }
    }

    /// Adds a (failed) requirement to the result collection. Adds an generic message to the collection
    /// if no on_fail message was specified.
    pub fn add<T: Fn() -> bool>(&mut self, r: &Requirement<T>) {
        self.push(match &r.on_fail {
            Some(ref x) => x.to_string(),
            None => "Some Requirement test failed but there is no explanation for it".to_string(),
        });
    }

    /// Adds an error message that some Requirement doesn't contain a check closure.
    pub fn log_no_check(&mut self) {
        self.push("Got Requirement without a test. Please investigate.".to_string());
    }

    /// Outputs the Results to the standard output.
    pub fn output(&self) {
        for rsl in &self.messages {
            println!("{}", rsl)
        }
    }

    /// Pushes a new message to the collection and prepends it with the name of the record.
    fn push(&mut self, msg: String) {
        self.messages.push(format!("{}: {}", self.record_name, msg));
    }
}

/// A Requirement represents a single check for one statement. Multiple Requirements make
/// a Check witch checks a specific Record. Each Requirement should only test one case.
pub struct Requirement<T>
where
    T: Fn() -> bool,
{
    /// Tests a certain requirement, returns true if this requirements is satisfied.
    check: Option<T>,
    /// The message shown to the user when the requirement doesn't hold. Should explain why the
    /// Record doesn't hold the requirement.
    on_fail: Option<String>,
}

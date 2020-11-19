use super::common::{Ident, ID};
use super::date::Date;
use super::entity::{Customer, Employee, Entity};
use super::money::Money;
use super::project::Project;
use super::record::{Record, RecordType};
use super::relation::Relation;
use super::transaction::Transaction;

use std::fmt;
use std::path::PathBuf;

use serde::{Deserialize, Serialize};

/// The Expense Categories assign each type of Expense with the appropriate expense-account in the
/// ledger. This is necessary as different types of expenses are treated differently in the
/// accounting. For example your travel expenses are booked into another account as the social
/// insurance bill.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExpenseCategory<'a> {
    /// Internal unique identifier of the Expense.
    id: ID,
    /// User-chosen and human readable identifier. This is helpful to mark a record and it's
    /// attachments more understandable as only using some long UUID.
    ident: Ident,
    /// Brief description of the purpose of the Expense Category.
    name: &'a str,
    /// HLedger conform account-path where the expense will added. The sub-accounts are
    /// separated by a colon (`:`). Example:
    ///
    /// ```
    /// 4 Aufwand:40 Materialaufwand:400 Materialaufwand
    /// ```
    account: &'a str,
}

impl<'a> Default for ExpenseCategory<'a> {
    fn default() -> Self {
        Self {
            id: ID::new(),
            ident: Ident::from_n(RecordType::ExpenseCategory, 1),
            name: "Costs of Materials",
            account: "4 Aufwand:40 Materialaufwand:400 Materialaufwand",
        }
    }
}

impl<'a> fmt::Display for ExpenseCategory<'a> {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "Expense Category {} ({})", self.name, self.ident)
    }
}

impl<'a> Record for ExpenseCategory<'a> {
    fn id(&self) -> ID {
        return ID::new();
    }
    fn ident(&self) -> Ident {
        return Ident::from("ec-23");
    }
    fn set_ident(&mut self, ident: Ident) {}
    fn record_type(&self) -> RecordType {
        RecordType::ExpenseCategory
    }
}

/// Categorizes the way a expense was paid. This is needed to create the correct transaction within
/// the ledger.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum PaymentMethod {
    /// The expense was paid with cash.
    Cash,
    /// The expense was paid using a credit card.
    Credit,
    /// The expense was paid with a debit card. Thus the amount was deducted directly from the bank
    /// account.
    Debit,
    /// Payed by bank-transfer (default).
    BankTransfer,
}

impl Default for PaymentMethod {
    fn default() -> Self {
        PaymentMethod::BankTransfer
    }
}

impl fmt::Display for PaymentMethod {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        f.write_str(match self {
            PaymentMethod::Cash => "Cash",
            PaymentMethod::Credit => "Credit",
            PaymentMethod::Debit => "Debit",
            PaymentMethod::BankTransfer => "Bank Transfer",
        })
    }
}

#[derive(Debug, Clone, Deserialize, Serialize)]
/// Expense represents a payment done by the company or a third party to assure the ongoing of
/// the business.
pub struct Expense<'a> {
    /// Internal unique identifier of the Expense.
    id: ID,
    /// User-chosen and human readable identifier. This is helpful to mark a record and it's
    /// attachments more understandable as only using some long UUID.
    ident: Ident,
    /// Describes the Expense in a meaningful way.
    name: &'a str,
    /// States the amount payed for the Expense.
    amount: Money,
    /// Path to an attachment which serves as a business record for this expense. The most common
    /// example is the scan of an invoice or a receipt.
    path: PathBuf,
    /// The day the obligation of this expense emerged.
    date_of_accrual: Date,
    /// States whether the Expense has to be forwarded to any third party.
    billable: bool,
    /// Refers to the customer which has to pay for the expense (if this is the case).
    obliged_customer: Option<Relation<Entity<'a, Customer>>>,
    /// States whether any third party (most common: a employee) has advanced this expense for the
    /// company and needs a pay back.
    advanced_by_third_party: bool,
    /// Refers to the third party which advanced the payment.
    advanced_third_party: Option<Relation<Entity<'a, Employee>>>,
    /// The date of the settlement of the expense (the company has not to take further actions).
    date_of_settlement: Option<Date>,
    /// The bank transaction which settled the Expense for the company if one exists.
    settlement_transaction: Option<Relation<Transaction>>,
    /// Categorizes the expense. This information is needed to create the correct transactions
    /// within the ledger.
    expense_category: Option<Relation<ExpenseCategory<'a>>>,
    /// The method used to pay for the expense. This information is needed to create the
    /// appropriate entries in the ledger.
    payment_method: PaymentMethod,
    /// True if the expense was for internal proposes only.
    internal: bool,
    /// References to the project for which the expense was paid (if this is the case).
    project: Option<Relation<Project>>,
}

impl<'a> Default for Expense<'a> {
    fn default() -> Self {
        Self {
            id: ID::new(),
            ident: Ident::from_n(RecordType::Expense, 1),
            name: "HAL 9000",
            amount: Money::default(),
            path: PathBuf::from("/path/to/attachement.pdf"),
            date_of_accrual: Date::default(),
            billable: false,
            obliged_customer: None,
            advanced_by_third_party: false,
            advanced_third_party: None,
            date_of_settlement: None,
            settlement_transaction: None,
            expense_category: None,
            payment_method: PaymentMethod::default(),
            internal: false,
            project: None,
        }
    }
}

impl<'a> fmt::Display for Expense<'a> {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(
            f,
            "Expense {} ({}) for {}",
            self.name, self.ident, self.amount
        )
    }
}

impl<'a> Record for Expense<'a> {
    fn id(&self) -> ID {
        return self.id;
    }
    fn ident(&self) -> Ident {
        return self.ident.clone();
    }
    fn set_ident(&mut self, ident: Ident) {}
    fn record_type(&self) -> RecordType {
        RecordType::Expense
    }
}

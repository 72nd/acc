mod schema;
mod test;

use schema::{Customer, Entity};

fn main() {
    let t: Entity<Customer> = Entity::default();
    println!("{:?}", t);
}

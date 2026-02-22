# Normalization Details

Normal forms progression with violation examples and resolutions.

---

## First Normal Form (1NF)

**Rule:** Eliminate repeating groups; each cell contains atomic values.

**Violation Example:**
```
Order(id, customer, items: "widget,gadget,thing")
```

**Resolution:**
```
Order(id, customer)
OrderItem(order_id, item_name)
```

## Second Normal Form (2NF)

**Rule:** Remove partial dependencies on composite keys.

**Violation Example:**
```
OrderItem(order_id, product_id, product_name, quantity)
                                 ^-- depends only on product_id
```

**Resolution:**
```
OrderItem(order_id, product_id, quantity)
Product(product_id, product_name)
```

## Third Normal Form (3NF)

**Rule:** Remove transitive dependencies; non-key columns depend only on the key.

**Violation Example:**
```
Employee(id, department_id, department_name)
                            ^-- depends on department_id, not employee id
```

**Resolution:**
```
Employee(id, department_id)
Department(id, name)
```

## Boyce-Codd Normal Form (BCNF)

**Rule:** Every determinant is a candidate key.

**Violation Example:**
```
CourseOffering(student, course, instructor)
-- Constraint: each instructor teaches only one course
-- instructor -> course (but instructor is not a candidate key)
```

**Resolution:**
```
InstructorCourse(instructor, course) -- instructor is key
Enrollment(student, instructor) -- references instructor
```

## When to Stop Normalizing

Stop at 3NF for most OLTP systems. Consider BCNF when:
- Update anomalies cause data corruption
- Data integrity is paramount
- Write frequency is high

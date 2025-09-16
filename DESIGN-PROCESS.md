## Understanding the framework

Laravel is a MVC fully contained ("complete ecosystem") web framework that utilizes php.

Comparison -> Php is to Laravel as Python is to Django

Laravel has an opinionated way of formatting projects such that it will handle the control of data, views and routing.

We don't need to set up the database ourselves (run migrations to get schemas set up.) We can do this via php artisan migrate where Laravel will automatically look for files in the database/migration folder within my project.

## Database Migrations

Facades are like Java helper classes that public static all of their methods within the class to be used throughout the project. Another analogy would be Django shortcut functions.

Common functionalities are used throughout the Laravel project for common necesities via these facades. With database migrations, the Schema facade is used for creating database schema. They require an up and down method (like most migrations)

## Routes

The routes folder... defines routes! Utilizes the Route facade for defining these routes. Methods are the http action words (GET, POST, PUT, PATCH, DELETE, OPTIONS).

First argument is the route's path, second argument is an array with the controller that the route utilizes and the second element is the method within that controller.

Likely, this was set up using Laravel Sanctum (project has the default route/api.php with token authentication that comes with running php artisan install:api)

Routes can handle redirects via the redirect or permenantRedirect methods from the Route facade.
Routes can handle serving views via the views method from the Route facade.

For this backend project, this won't be necessary.

## Views

Laravel uses blade templates for serving html to users. Found in resource/views. Won't go any further since this is mainly a backend api application.

## Controllers

Controllers bundle together the common functions for dealing with a certain model as they do in most MVC architecture.

Controllers get access to requests via Illuminate\Http\Request

## Models

Models map data typically fairly 1:1 with a database's schema (again like most MVCs). These are the entry pont for understanding Laravel's "Eloquent ORM".

Models automatically map to a table via snake case. If the model name doesn't exactly match the table name, we can set the protected table property of a model for Laravel to pick it up manually.

I wanted to create a reminders model with php artisan make:mdoel Reminders, but I don't have php installed on my laptop so I'll just make them manually

## Resources

Laravel's Resource API is how models are transformed to be served to users. This can be seen in the UserController as the GET endpoint is returning a UserResource

What confused me when initially reading was the toArray method in the sample UserResource class. Just from the name, I'd expect the response to be an array. However, it is in the format {"data": {"id": 1, "email"...}} much more formatted and not a direct array of values.

However it seems that the toArray method is Laravel's way of defining the mapping of properties to values when returning responses.

Got me curious about how returning multiple data objects would look like. Luckily, it seems that Laravel an opinionated way for handling this as well with the collections method on a Resource's class.

_NOTE: Laravel's documentation is really well made and the flow of understanding concepts is very smooth_

Laravel even handles transformation of a Model into a Resource with `toResource` and `toResourceCollection`

## Linking back to User sample

1. What is findOrFail
   Laravel's eloquent ORM way of fetching a single resource from the database by its primary key,and returns throws the user a 404 error if it can't find it.

2. What would the data flow look like from request to response
   From the request, Laravel maps it to the specific `route` which sends it to the `controller` and its corresponding `method`. Using Eloquent ORM, it fetches the data via the `model`. It then goes through a transformation layer via the `resource` since the controller is returning the UserResource which only returns specific fields that we want the api user to see (i.e we don't want to be returning the user's password or remember_token). `toArray` defines how Laravel will map the data. Laravel will wrap this in a data key and convert it to json to be returned in the response.

    So request --> controller --> model --> resource --> response

# Recurring Reminders

Minimum key functional requirements:

1. Create, update, delete reminders and their recurrence rules
2. View a list of reminders that will occur within a given date range, based on their recurrence rules
3. Searching for reminders based on a keyword
4. Allow for multiple kinds of recurrence rules ("every day", "every n days/weeks/months", "every second Monday")

My first instinct was to have some sort of equation such that we can calculate the all the reminders needed. If we have the start and end date and the rule, then we can determine all the reminders inbetween. Although this runs into an issue where if a user doesn't want to define an end date (like "every day", "every 2 weeks") then this would result in an infinite number of reminders. There must be some way to calculate the dates without having to rely on an end date. Quick searching lead me to iCalendar's recurrence rule standard (RFC 5545). With this we can easily calculate any date.

## RFC 5545 3.3.10 - Reccurence Rules (RRULE)

List of key value pairs that define how something repeats over time. Requires FREQ and the other keys are optional.

FREQ -> "SECONDLY" / "MINUTELY" / "HOURLY" / "DAILY" / "WEEKLY" / "MONTHLY" / "YEARLY" (base unit of repetition)

UNTIL -> date / date-time (date[time] it should end)

COUNT -> int (how many times to repeat)

INTERVAL -> intervals of repetition

BY\* -> which occurence to repeat

\*NOTE: UNIT AND COUNT ARE MUTUALLY EXCLUSIVE

With this, we can easily determine which dates for a reminder are within a certain range.

# Design

## Minimum Functional Requirements

To address the minimum functional requirements:

1. Create, update, delete reminders and their recurrence rules
    - POST /api/reminder
    - PATCH/PUT /api/reminder
    - DELETE /api/reminder
2. View a list of reminders that will occur within a given date range, based on their recurrence rules
    - GET /api/reminder?startDate={startDate}?endDate={endDate}
3. Searching for reminders based on a keyword
    - Domain should have some description to search on
4. Allow for multiple kinds of recurrence rules ("every day", "every n days/weeks/months", "every second Monday")
    - Utilize RFC 5545

## Database Schema

We can create a database schema such that it addressed the above:

```
CREATE TABLE reminders (
    id SERIAL PRIMARY KEY
    user_id INT REFERENCES users(id)
    rrule TEXT
    description TEXT
    start_date DATETIME
)
```

Since recurrence rules already define how they expire, we wouldn't need to store the end_date directly in our table.

## Model

The model for this domain object follows easily. We can define a migration and run php artisan make:model Reminder

## Resource

How we want to return the recurrence rule is the challenging aspect. We'd have to define whether returning the recurrence rule as is (in RFC 5545 format), is sufficient for the requirements. For MVP, we can add a route GET /api/reminder/{reminderId} that simply returns the resource with the recurrence rule as is.

## Controller

Starting with the first three routes, these should follow fairly simply. Running make:model earlier should have created this for us.

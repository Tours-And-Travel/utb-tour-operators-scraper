I want you to extract structured data from the following web page:
URL: https://www.hiddengemsofuganda.com/gorilla-tours

Your task is to:

Parse all the tour packages listed on the page.

For each tour, extract and structure the following fields:

name (string): The name/title of the tour

description (string): A brief summary of what the tour includes

duration_days (number): How many days the tour lasts

starting_location (array of strings): All listed starting points

regions (array): Areas covered in the tour (e.g., Bwindi, Queen Elizabeth NP)

images (array of strings): The images representing the tour

type (string): "Private" or "Group" tour

min_group_size (number): Minimum group size

includes (array of strings): List of services included (e.g., permits, meals, transport)

excludes (array of strings): List of services not included

gorilla_permit_included (boolean): Whether gorilla permit is included

accommodation (object):

type (string): e.g. Lodge, Hotel

included (boolean): whether it's included

availability (array): Dates or general availability like "Daily" or "On Request"

language (array): Languages spoken by guides

difficulty_level (string): e.g. Easy / Moderate / Hard

transport_type (string): e.g. 4x4, van

customizable (boolean): Whether the tour can be customized

booking_url (string): A link to book or contact

prices (object):

solo (number): Price for solo traveller

group (object): Keys as group sizes (e.g. "2", "3", "4+"), values as per-person prices

currency (string): e.g. "USD"

tour_id (string): A unique tour code (you can generate this) e.g. "UG-GT-001"

Output a single JSON object with an "activities" array containing all the tour objects.

Ensure all URLs are absolute.

Ensure accurate prices from the website please

Clean the data: remove extra line breaks or HTML tags, and trim whitespace.

❗ Do not invent content. Only include fields that are clearly visible or can be logically inferred (e.g., if it says “starts from Kigali,” include "Kigali Airport" in starting_location).


✅ Format your output as valid JSON only, no extra text.

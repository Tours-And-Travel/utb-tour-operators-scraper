I want you to extract structured travel itinerary data from the following web page:
URL: https://labaafrica.com/tour/8-days-uganda-gorillas-chimpanzee-and-tree-climbing-lions-tour/

Your task is to:

Parse the daily itinerary schedule, capturing information for each day of the tour.

For each day, extract and structure the following fields:

- **day** (number): The day number
- **title** (string): Title or name of the day’s experience
- **location** (string): The location(s) covered that day
- **images** (array of strings): The images representing the tour (you can find this below the itinerary item details)
- **coordinates** (array of numbers): GPS coordinates in [latitude, longitude]. (you can generate this, use a tool to lookup coordinates)
- **accommodation** (string): Name of the accommodation for the night, or "End of Safari" if it's the last day
- **highlights** (array of strings): Key attractions or experiences of the day
- **activities** (array of objects):
  - **icon** (string): Symbol representing the activity (you can generate this based on https://lucide.dev) e.g airplane
  - **title** (string): Activity title
  - **time** (string): Time of the activity (e.g., "Morning", "Afternoon")
  - **description** (string): Brief summary of the activity

Format the entire itinerary as a JSON object with a single top-level field:

```json
{
  "itinerary": [/* each day's structured data here */]
}

✅ Format your output as valid JSON only, no extra text.

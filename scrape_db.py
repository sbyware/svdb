import requests
from bs4 import BeautifulSoup
import json

def scrape_wikipedia(url):
    response = requests.get(url)
    if response.status_code == 200:
        soup = BeautifulSoup(response.content, 'html.parser')
        # Find the table containing port numbers
        table = soup.find('table', class_='wikitable')
        if table:
            data = []
            headers = [header.text.strip() for header in table.find_all('th')]
            for row in table.find_all('tr')[1:]:
                cells = row.find_all(['td', 'th'])
                entry = {}
                for i, cell in enumerate(cells):
                    if headers[i] == 'Port':
                        entry['Port'] = int(cell.text.strip())
                    elif headers[i] == 'TCP':
                        entry['TCP'] = cell.text.strip() if cell.text.strip() != 'No' else None
                    elif headers[i] == 'UDP':
                        entry['UDP'] = cell.text.strip() if cell.text.strip() != 'No' else None
                data.append(entry)
            return data
        else:
            print("Table not found on the page.")
            return None
    else:
        print("Failed to retrieve the Wikipedia page.")
        return None

def save_to_json(data, filename):
    with open(filename, 'w') as f:
        json.dump(data, f, indent=4)

def main():
    wikipedia_url = "https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers"
    data = scrape_wikipedia(wikipedia_url)
    if data:
        save_to_json(data, "svdb-scraped.json")
        print("Data saved to svdb.json successfully.")
    else:
        print("Failed to scrape data.")

if __name__ == "__main__":
    main()

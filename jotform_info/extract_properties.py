from bs4 import BeautifulSoup
from glob import glob
import json

# usage: python extract_properties.py --pivot *.html > jotform_properties.json

def convert(files):
    data = []
    # files=glob("*.html", recursive=False)
    for file in files:
        with open(file) as f:
            soup = BeautifulSoup(f.read(), features="html.parser")
        table = soup.findAll("table")[1]
        keys = [th.string for th in table.findAll("th")]
        data = data + [dict(zip(keys,values)) for values in [[td.string for td in tr.findAll("td")] for tr in table.findAll("tr")[1:]]]

    return data

def field_dict(data):
    data_dict = {}
    for prop in data:
        if not data_dict.get(prop["Field"]):
            data_dict[prop["Field"]] = [prop]
        else:
            data_dict[prop["Field"]].append(prop)
    return data_dict

def prop_dict(data):
    """
    pivot to propname: (list of fields)
    :param data:
    :return:
    """
    data_dict = {}
    for prop in data:
        if not data_dict.get(prop["Property"]):
            data_dict[prop["Property"]] = [prop]
        else:
            data_dict[prop["Property"]].append(prop)

if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("file", nargs="+")
    parser.add_argument("--pivot", action="store_true", default=False, help="Pivot into a dictionary where Field values are keys")
    args = parser.parse_args()

    data = convert(args.file)
    if args.pivot:
        data_dict = field_dict(data)
        print(json.dumps(data_dict, indent=2))
    else:
        print("\n".join([json.dumps(d) for d in data]))


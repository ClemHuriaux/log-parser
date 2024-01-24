import pandas as pd
import glob
import sys
import argparse
import time


def read_file(file_path: str):
    with open(file_path, 'r', encoding="utf8") as file:
        return file.read().splitlines()
       

def parse_logs(source, destination, columns_file, verbose=False, sep=' '):
    all_files = glob.glob(fr"{source}\*.*")
    columns = read_file(columns_file)
    all_files_len = len(all_files)
   
    data = []
    for i, files in enumerate(all_files, start=1):
        if verbose:
            print(f"Parsing file {i}/{all_files_len}")
            
        file = read_file(files)
        data.extend([line.split(sep)[:len(columns)] for line in file])
            
        if verbose:
            progress_percentage = (i / all_files_len) * 100
            print(f"Progress: {progress_percentage:.2f}%")
        
    df = pd.DataFrame(data, columns=columns)
    print("Saving to csv")
    df.to_csv(f'{destination}\\results.csv', index=False)
    print("Done!")
    

if __name__ == "__main__":
    start_time = time.time()
    parser = argparse.ArgumentParser(description="This tool is designed to parse logfiles having the structure of a txtfile",
                                     formatter_class=argparse.ArgumentDefaultsHelpFormatter,
                                     epilog="Example: python log_parser.py C:\columns.txt C:\destination C:\source -v -s '-' ")
    parser.add_argument("-v", "--verbose", action="store_true", default=False, help="Show state of files parsing")
    parser.add_argument("-s", "--separator", default=' ', help="Specify the separator of the files")
    parser.add_argument("columns", help="Provide the full path of the file containing the columns name. DON'T forget the extension")
    parser.add_argument("destination", help="Give the full path of the folder where the csv will be saved")
    parser.add_argument("source", help="Give the full path of the folder containing the log files")
    args = parser.parse_args()
    parse_logs(args.source, args.destination, args.columns, args.verbose, args.separator)
    end_time = time.time()
    print(f'Execution time: {end_time - start_time} seconds')
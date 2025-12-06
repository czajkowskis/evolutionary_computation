import pandas as pd
import matplotlib.pyplot as plt
import glob

def plot_results(filename):
    try:
        df = pd.read_csv(filename)
    except FileNotFoundError:
        print(f"File {filename} not found.")
        return

    instance_name = filename.split('_')[-1].replace('.csv', '')
    
    # Define the 6 plots
    plots = [
        ('Objective', 'SimAvgEdges', f'Instance {instance_name}: Objective vs Avg Similarity (Edges)'),
        ('Objective', 'SimAvgNodes', f'Instance {instance_name}: Objective vs Avg Similarity (Nodes)'),
        ('Objective', 'SimBest1000Edges', f'Instance {instance_name}: Objective vs Similarity to Best of 1000 (Edges)'),
        ('Objective', 'SimBest1000Nodes', f'Instance {instance_name}: Objective vs Similarity to Best of 1000 (Nodes)'),
        ('Objective', 'SimBestKnownEdges', f'Instance {instance_name}: Objective vs Similarity to Best Known (Edges)'),
        ('Objective', 'SimBestKnownNodes', f'Instance {instance_name}: Objective vs Similarity to Best Known (Nodes)'),
    ]

    # Get reference objectives
    best_known_obj = df['BestKnownObjective'].iloc[0] if 'BestKnownObjective' in df.columns else None
    
    # Best of 1000 is the minimum objective in the dataset (or specifically marked)
    # The dataset contains 1000 solutions. The best one is the min.
    best_of_1000_obj = df['Objective'].min()
    
    # Average objective
    avg_obj = df['Objective'].mean()

    for x_col, y_col, title in plots:
        plt.figure(figsize=(10, 6))
        
        # Filter data if needed (e.g., exclude best solution for Best1000 plots)
        data = df
        if 'Best1000' in y_col:
            data = df[df['IsBestOf1000'] == False]
        
        plt.scatter(data[x_col], data[y_col], alpha=0.5, label='Local Optima')
        
        # Add reference lines
        if 'BestKnown' in y_col and best_known_obj is not None:
            plt.axvline(x=best_known_obj, color='r', linestyle='--', label=f'Best Known ({best_known_obj})')
        elif 'Best1000' in y_col:
            plt.axvline(x=best_of_1000_obj, color='r', linestyle='--', label=f'Best of 1000 ({best_of_1000_obj})')
        elif 'SimAvg' in y_col:
            plt.axvline(x=avg_obj, color='r', linestyle='--', label=f'Average Objective ({avg_obj:.2f})')

        plt.title(title)
        plt.xlabel('Objective Function Value')
        plt.ylabel('Similarity')
        plt.legend(loc='lower right')
        plt.grid(True)
        
        # Dynamic y-axis limits: Min 10 (or 70 for Nodes), Max = max value + 5
        max_val = data[y_col].max()
        min_ylim = 70 if 'Nodes' in y_col else 10
        plt.ylim(min_ylim, max_val + 5)



        
        # Save plot
        safe_title = title.replace(' ', '_').replace(':', '').replace('(', '').replace(')', '')
        plt.savefig(f'{safe_title}.png')
        plt.close()
        print(f"Saved {safe_title}.png")

def main():
    files = glob.glob("convexity_results_*.csv")
    for f in files:
        plot_results(f)

if __name__ == "__main__":
    main()

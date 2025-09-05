import time

def main():
    """
    This function runs a speed test loop in Python.
    """
    # The number of iterations for the loop.
    # You can change this value to test different loop sizes.
    iterations = 1_000_000_000

    # Start a timer to measure the execution time.
    start_time = time.time()

    # Initialize a counter. We'll perform a simple operation inside the loop.
    counter = 0

    # The loop for the speed test.
    for i in range(iterations):
        counter = i

    # Calculate the duration of the loop.
    duration = time.time() - start_time

    # Print the results.
    print(f"Python loop with {iterations:,} iterations took: {duration:.4f} seconds")
    # The final value of the counter will be iterations - 1.
    print(f"Final counter value: {counter}")

if __name__ == "__main__":
    main()

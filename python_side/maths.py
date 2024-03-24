from operator import itemgetter

# start partition problem algorithm from https://stackoverflow.com/a/7942946
# modified to act on list of images rather than the weights themselves
# more info on the partition problem http://www8.cs.umu.se/kurser/TDBAfl/VT06/algorithms/BOOK/BOOK2/NODE45.HTM


def linear_partition(sequence, num_rows, data_list):
    min_l = len(sequence) - 1
    if num_rows > min_l:
        return map(lambda x: [x], sequence)

    solution = linear_partition_table(sequence, num_rows)
    num_rows = num_rows - 2
    answer = []

    while num_rows >= 0:
        answer = [
            [data_list[i] for i in range(solution[min_l - 1][num_rows] + 1, min_l + 1)]
        ] + answer

        min_l = solution[min_l - 1][num_rows]
        num_rows = num_rows - 1

    answer = [[data_list[i] for i in range(0, min_l + 1)]] + answer

    # print(f"linear_partition({ans=})")
    return answer


def linear_partition_table(sequence, num_rows):
    # print(f"linear_partition_table({sequence=}, {num_rows=})")
    num_elements = len(sequence)
    table = []
    solution = []

    for _ in range(num_elements):
        table.append([0] * num_rows)

    for _ in range(num_elements - 1):
        solution.append([0] * (num_rows - 1))

    for index in range(num_elements):
        table[index][0] = sequence[index] + (table[index - 1][0] if index else 0)

    for col_idx in range(num_rows):
        table[0][col_idx] = sequence[0]

    for index in range(1, num_elements):
        for col_idx in range(1, num_rows):
            optimal_partition = []

            for x in range(index):
                max_value = max(table[x][col_idx - 1], table[index][0] - table[x][0])
                tuple_element = (max_value, x)
                optimal_partition.append(tuple_element)

            table[index][col_idx], solution[index - 1][col_idx - 1] = min(
                optimal_partition,
                key=itemgetter(0),
            )
    # print(f"table_solution: {solution}")
    return solution


# end partition problem algorithm


def clamp(v, h):
    return 1 if v < 1 else h if v > h else v


def ensure_even(n):
    return n if n % 2 == 0 else n + 1

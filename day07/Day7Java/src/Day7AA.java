import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;

public class Day7AA {

    private static Pair<Long, List<Long>> parseLine(String line) {
        String[] parts = line.split(": ");
        long target = Long.parseLong(parts[0]);
        String[] numbersStr = parts[1].split(" ");
        List<Long> numbers = new ArrayList<>();
        for (String numStr : numbersStr) {
            numbers.add(Long.parseLong(numStr));
        }
        return new Pair<>(target, numbers);
    }

    private static boolean isValid(long target, List<Long> values, boolean allowConcat) {
        if (values.size() == 1) {
            return values.getFirst() == target;
        }

        List<Long> additionList = new ArrayList<>(values);
        additionList.set(0, values.get(0) + values.get(1));
        additionList.remove(1);
        if (isValid(target, additionList, allowConcat)) {
            return true;
        }

        List<Long> multiplicationList = new ArrayList<>(values);
        multiplicationList.set(0, values.get(0) * values.get(1));
        multiplicationList.remove(1);
        if (isValid(target, multiplicationList, allowConcat)) {
            return true;
        }

        if (allowConcat) {
            String concatStr = values.get(0) + "" + values.get(1);
            long concatVal = Long.parseLong(concatStr);
            List<Long> concatList = new ArrayList<>(values);
            concatList.set(0, concatVal);
            concatList.remove(1);
            return isValid(target, concatList, true);
        }

        return false;
    }

    private static long part1(List<String> input) {
        long total = 0;
        for (String line : input) {
            Pair<Long, List<Long>> parsedLine = parseLine(line);
            long target = parsedLine.first();
            List<Long> numbers = parsedLine.second();
            if (isValid(target, numbers, false)) {
                total += target;
            }
        }
        return total;
    }

    private static long part2(List<String> input) {
        long total = 0;
        for (String line : input) {
            Pair<Long, List<Long>> parsedLine = parseLine(line);
            long target = parsedLine.first();
            List<Long> numbers = parsedLine.second();
            if (isValid(target, numbers, true)) {
                total += target;
            }
        }
        return total;
    }

    public static void main(String[] args) {
        try {
            List<String> input = Files.readAllLines(Path.of("../input.txt"));

            System.out.println("Part 1: " + part1(input));
            System.out.println("Part 2: " + part2(input));
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}

// Utility Pair class for returning two values
record Pair<F, S>(F first, S second) {
}

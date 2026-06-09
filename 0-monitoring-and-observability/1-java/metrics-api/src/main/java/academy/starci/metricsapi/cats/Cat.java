package academy.starci.metricsapi.cats;

/**
 * Demo resource returned by GET /cats and created by POST /cats.
 */
public record Cat(int id, String name, int age) {
}

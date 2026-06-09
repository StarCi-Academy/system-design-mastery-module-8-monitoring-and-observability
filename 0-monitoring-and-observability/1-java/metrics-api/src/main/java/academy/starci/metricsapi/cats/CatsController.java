package academy.starci.metricsapi.cats;

import jakarta.validation.Valid;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.atomic.AtomicInteger;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

/**
 * REST controller for /cats. In-memory store seeded with one row (Tom) so the
 * first GET is meaningful and Prometheus sees real traffic.
 */
@RestController
@RequestMapping("/cats")
public class CatsController {

    private final List<Cat> cats = new CopyOnWriteArrayList<>(List.of(new Cat(1, "Tom", 3)));
    private final AtomicInteger nextId = new AtomicInteger(2);

    @GetMapping
    public List<Cat> findAll() {
        return cats;
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public Cat create(@Valid @RequestBody CreateCatRequest body) {
        Cat cat = new Cat(nextId.getAndIncrement(), body.getName(), body.getAge());
        cats.add(cat);
        return cat;
    }
}

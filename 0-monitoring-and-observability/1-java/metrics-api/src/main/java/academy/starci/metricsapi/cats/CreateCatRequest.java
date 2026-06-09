package academy.starci.metricsapi.cats;

import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;

/**
 * Request body for POST /cats. name is required (NotEmpty), age must be present.
 */
public class CreateCatRequest {

    @NotEmpty(message = "name should not be empty")
    private String name;

    @NotNull(message = "age must be an integer number")
    private Integer age;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public Integer getAge() {
        return age;
    }

    public void setAge(Integer age) {
        this.age = age;
    }
}

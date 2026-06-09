package academy.starci.metricsapi.cats;

import java.util.List;
import java.util.Map;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

/**
 * Shapes validation failures into the shared error contract:
 * { "statusCode": 400, "message": ["..."], "error": "Bad Request" }.
 */
@RestControllerAdvice
public class ValidationExceptionHandler {

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<Map<String, Object>> handleValidation(MethodArgumentNotValidException ex) {
        List<String> messages = ex.getBindingResult().getFieldErrors().stream()
                .map(error -> error.getDefaultMessage())
                .toList();
        Map<String, Object> body = Map.of(
                "statusCode", 400,
                "message", messages,
                "error", "Bad Request");
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(body);
    }
}

package io.respondnow.util.constants;

public final class AppConstants {
    private AppConstants() {
        throw new UnsupportedOperationException("Constants class cannot be instantiated.");
    }


    /**
     * API Endpoint Paths
     */
    public static final class ApiPaths {
        public static final String AUTH_BASE = "/auth";
        public static final String SIGNUP = "/signup";
        public static final String LOGIN = "/login";
        public static final String CHANGE_PASSWORD = "/changePassword";
        public static final String USER_MAPPING = "/userMapping";

        private ApiPaths() {
            throw new UnsupportedOperationException("ApiPaths constants class cannot be instantiated.");
        }
    }

    public static final class ResponseStatus {
        public static final String SUCCESS = "success";
        public static final String ERROR = "error";
        
        private ResponseStatus() {
            throw new UnsupportedOperationException("ResponseStatus constants class cannot be instantiated.");
        }
    }

    public static final class ErrorMessages {

        private ErrorMessages() {
            throw new UnsupportedOperationException("ApiPaths constants class cannot be instantiated.");
        }
    }


}

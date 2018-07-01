package com.sidekick.dataservice.component;

import org.springframework.stereotype.Component;

import java.util.HashMap;
import java.util.Map;

@Component
public class ResultMap {

    public Map<String, Object> successMap(Object obj) {
        HashMap<String, Object> mapObj = new HashMap<>();
        mapObj.put("success", true);
        mapObj.put("message", "");
        mapObj.put("code", 0);

        if (obj != null) {
            mapObj.put("data", obj);
        }

        return mapObj;
    }

    public Map<String, Object> success() {
        return successMap(null);
    }

    public Map<String, Object> failedMap(Object reqobj, String message, Integer code) {
        HashMap<String, Object> mapObj = new HashMap<>();
        mapObj.put("success", false);

        if (message != null) {
            mapObj.put("message", message);
        }

        if (code != null) {
            mapObj.put("code", code);
        }

        if (reqobj != null) {
            mapObj.put("req", reqobj);
        }
        return mapObj;
    }

}

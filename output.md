# Keyrotator Package: Outputs and Errors

---

## **Data Return Types**
- **Key IDs:** `string`
- **Key Values:** `interface{}` (JSON-compatible types)
- **Errors:** `error` interface



### **NewAPIKeyConfig Errors**
1. **Error:** `"failed to get absolute path: [error details]"`
2. **Error:** `"failed to read config file: [error details]"`
3. **Error:** `"failed to parse config: [error details]"`

---

### **GetRandomAPIKey Outputs**
#### Success:  
- **Returns:** `(keyID string, keyValue interface{}, nil)`

#### Errors:  
1. **Error:** `"category [name] not found"`
2. **Error:** `"no available keys in category [name]"`

---

### **GetAPIKey Outputs**
#### Success:  
- **Returns:** `(keyValue interface{}, nil)`

#### Errors:  
1. **Error:** `"category [name] not found"`
2. **Error:** `"key [id] not found in category [name]"`

---

### **MarkKeyAsExhausted**
#### Success:  
- **Returns:** `nil` 

#### Errors:  
1. **Error:** `"category [name] not found"`
2. **Error:** `"key [id] not found in category [name]"`

---

### **Reload**
#### Success:  
- **Returns:** `nil`

#### Errors:  
- Any error from **NewAPIKeyConfig** can propagate.

---

### **SaveConfig**
#### Success:  
- **Returns:** `nil`

#### Errors:  
1. **Error:** `"failed to marshal config: [error details]"`
2. **Error:** `"failed to write config file: [error details]"`



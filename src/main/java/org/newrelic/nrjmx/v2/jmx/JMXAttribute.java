/**
 * Autogenerated by Thrift Compiler (0.14.1)
 *
 * DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING
 *  @generated
 */
package org.newrelic.nrjmx.v2.jmx;

@SuppressWarnings({"cast", "rawtypes", "serial", "unchecked", "unused"})
@javax.annotation.Generated(value = "Autogenerated by Thrift Compiler (0.14.1)", date = "2021-10-27")
public class JMXAttribute implements org.apache.thrift.TBase<JMXAttribute, JMXAttribute._Fields>, java.io.Serializable, Cloneable, Comparable<JMXAttribute> {
  private static final org.apache.thrift.protocol.TStruct STRUCT_DESC = new org.apache.thrift.protocol.TStruct("JMXAttribute");

  private static final org.apache.thrift.protocol.TField ATTRIBUTE_FIELD_DESC = new org.apache.thrift.protocol.TField("attribute", org.apache.thrift.protocol.TType.STRING, (short)1);
  private static final org.apache.thrift.protocol.TField VALUE_FIELD_DESC = new org.apache.thrift.protocol.TField("value", org.apache.thrift.protocol.TType.STRUCT, (short)2);

  private static final org.apache.thrift.scheme.SchemeFactory STANDARD_SCHEME_FACTORY = new JMXAttributeStandardSchemeFactory();
  private static final org.apache.thrift.scheme.SchemeFactory TUPLE_SCHEME_FACTORY = new JMXAttributeTupleSchemeFactory();

  public @org.apache.thrift.annotation.Nullable java.lang.String attribute; // required
  public @org.apache.thrift.annotation.Nullable JMXAttributeValue value; // required

  /** The set of fields this struct contains, along with convenience methods for finding and manipulating them. */
  public enum _Fields implements org.apache.thrift.TFieldIdEnum {
    ATTRIBUTE((short)1, "attribute"),
    VALUE((short)2, "value");

    private static final java.util.Map<java.lang.String, _Fields> byName = new java.util.HashMap<java.lang.String, _Fields>();

    static {
      for (_Fields field : java.util.EnumSet.allOf(_Fields.class)) {
        byName.put(field.getFieldName(), field);
      }
    }

    /**
     * Find the _Fields constant that matches fieldId, or null if its not found.
     */
    @org.apache.thrift.annotation.Nullable
    public static _Fields findByThriftId(int fieldId) {
      switch(fieldId) {
        case 1: // ATTRIBUTE
          return ATTRIBUTE;
        case 2: // VALUE
          return VALUE;
        default:
          return null;
      }
    }

    /**
     * Find the _Fields constant that matches fieldId, throwing an exception
     * if it is not found.
     */
    public static _Fields findByThriftIdOrThrow(int fieldId) {
      _Fields fields = findByThriftId(fieldId);
      if (fields == null) throw new java.lang.IllegalArgumentException("Field " + fieldId + " doesn't exist!");
      return fields;
    }

    /**
     * Find the _Fields constant that matches name, or null if its not found.
     */
    @org.apache.thrift.annotation.Nullable
    public static _Fields findByName(java.lang.String name) {
      return byName.get(name);
    }

    private final short _thriftId;
    private final java.lang.String _fieldName;

    _Fields(short thriftId, java.lang.String fieldName) {
      _thriftId = thriftId;
      _fieldName = fieldName;
    }

    public short getThriftFieldId() {
      return _thriftId;
    }

    public java.lang.String getFieldName() {
      return _fieldName;
    }
  }

  // isset id assignments
  public static final java.util.Map<_Fields, org.apache.thrift.meta_data.FieldMetaData> metaDataMap;
  static {
    java.util.Map<_Fields, org.apache.thrift.meta_data.FieldMetaData> tmpMap = new java.util.EnumMap<_Fields, org.apache.thrift.meta_data.FieldMetaData>(_Fields.class);
    tmpMap.put(_Fields.ATTRIBUTE, new org.apache.thrift.meta_data.FieldMetaData("attribute", org.apache.thrift.TFieldRequirementType.DEFAULT, 
        new org.apache.thrift.meta_data.FieldValueMetaData(org.apache.thrift.protocol.TType.STRING)));
    tmpMap.put(_Fields.VALUE, new org.apache.thrift.meta_data.FieldMetaData("value", org.apache.thrift.TFieldRequirementType.DEFAULT, 
        new org.apache.thrift.meta_data.StructMetaData(org.apache.thrift.protocol.TType.STRUCT, JMXAttributeValue.class)));
    metaDataMap = java.util.Collections.unmodifiableMap(tmpMap);
    org.apache.thrift.meta_data.FieldMetaData.addStructMetaDataMap(JMXAttribute.class, metaDataMap);
  }

  public JMXAttribute() {
  }

  public JMXAttribute(
    java.lang.String attribute,
    JMXAttributeValue value)
  {
    this();
    this.attribute = attribute;
    this.value = value;
  }

  /**
   * Performs a deep copy on <i>other</i>.
   */
  public JMXAttribute(JMXAttribute other) {
    if (other.isSetAttribute()) {
      this.attribute = other.attribute;
    }
    if (other.isSetValue()) {
      this.value = new JMXAttributeValue(other.value);
    }
  }

  public JMXAttribute deepCopy() {
    return new JMXAttribute(this);
  }

  @Override
  public void clear() {
    this.attribute = null;
    this.value = null;
  }

  @org.apache.thrift.annotation.Nullable
  public java.lang.String getAttribute() {
    return this.attribute;
  }

  public JMXAttribute setAttribute(@org.apache.thrift.annotation.Nullable java.lang.String attribute) {
    this.attribute = attribute;
    return this;
  }

  public void unsetAttribute() {
    this.attribute = null;
  }

  /** Returns true if field attribute is set (has been assigned a value) and false otherwise */
  public boolean isSetAttribute() {
    return this.attribute != null;
  }

  public void setAttributeIsSet(boolean value) {
    if (!value) {
      this.attribute = null;
    }
  }

  @org.apache.thrift.annotation.Nullable
  public JMXAttributeValue getValue() {
    return this.value;
  }

  public JMXAttribute setValue(@org.apache.thrift.annotation.Nullable JMXAttributeValue value) {
    this.value = value;
    return this;
  }

  public void unsetValue() {
    this.value = null;
  }

  /** Returns true if field value is set (has been assigned a value) and false otherwise */
  public boolean isSetValue() {
    return this.value != null;
  }

  public void setValueIsSet(boolean value) {
    if (!value) {
      this.value = null;
    }
  }

  public void setFieldValue(_Fields field, @org.apache.thrift.annotation.Nullable java.lang.Object value) {
    switch (field) {
    case ATTRIBUTE:
      if (value == null) {
        unsetAttribute();
      } else {
        setAttribute((java.lang.String)value);
      }
      break;

    case VALUE:
      if (value == null) {
        unsetValue();
      } else {
        setValue((JMXAttributeValue)value);
      }
      break;

    }
  }

  @org.apache.thrift.annotation.Nullable
  public java.lang.Object getFieldValue(_Fields field) {
    switch (field) {
    case ATTRIBUTE:
      return getAttribute();

    case VALUE:
      return getValue();

    }
    throw new java.lang.IllegalStateException();
  }

  /** Returns true if field corresponding to fieldID is set (has been assigned a value) and false otherwise */
  public boolean isSet(_Fields field) {
    if (field == null) {
      throw new java.lang.IllegalArgumentException();
    }

    switch (field) {
    case ATTRIBUTE:
      return isSetAttribute();
    case VALUE:
      return isSetValue();
    }
    throw new java.lang.IllegalStateException();
  }

  @Override
  public boolean equals(java.lang.Object that) {
    if (that instanceof JMXAttribute)
      return this.equals((JMXAttribute)that);
    return false;
  }

  public boolean equals(JMXAttribute that) {
    if (that == null)
      return false;
    if (this == that)
      return true;

    boolean this_present_attribute = true && this.isSetAttribute();
    boolean that_present_attribute = true && that.isSetAttribute();
    if (this_present_attribute || that_present_attribute) {
      if (!(this_present_attribute && that_present_attribute))
        return false;
      if (!this.attribute.equals(that.attribute))
        return false;
    }

    boolean this_present_value = true && this.isSetValue();
    boolean that_present_value = true && that.isSetValue();
    if (this_present_value || that_present_value) {
      if (!(this_present_value && that_present_value))
        return false;
      if (!this.value.equals(that.value))
        return false;
    }

    return true;
  }

  @Override
  public int hashCode() {
    int hashCode = 1;

    hashCode = hashCode * 8191 + ((isSetAttribute()) ? 131071 : 524287);
    if (isSetAttribute())
      hashCode = hashCode * 8191 + attribute.hashCode();

    hashCode = hashCode * 8191 + ((isSetValue()) ? 131071 : 524287);
    if (isSetValue())
      hashCode = hashCode * 8191 + value.hashCode();

    return hashCode;
  }

  @Override
  public int compareTo(JMXAttribute other) {
    if (!getClass().equals(other.getClass())) {
      return getClass().getName().compareTo(other.getClass().getName());
    }

    int lastComparison = 0;

    lastComparison = java.lang.Boolean.compare(isSetAttribute(), other.isSetAttribute());
    if (lastComparison != 0) {
      return lastComparison;
    }
    if (isSetAttribute()) {
      lastComparison = org.apache.thrift.TBaseHelper.compareTo(this.attribute, other.attribute);
      if (lastComparison != 0) {
        return lastComparison;
      }
    }
    lastComparison = java.lang.Boolean.compare(isSetValue(), other.isSetValue());
    if (lastComparison != 0) {
      return lastComparison;
    }
    if (isSetValue()) {
      lastComparison = org.apache.thrift.TBaseHelper.compareTo(this.value, other.value);
      if (lastComparison != 0) {
        return lastComparison;
      }
    }
    return 0;
  }

  @org.apache.thrift.annotation.Nullable
  public _Fields fieldForId(int fieldId) {
    return _Fields.findByThriftId(fieldId);
  }

  public void read(org.apache.thrift.protocol.TProtocol iprot) throws org.apache.thrift.TException {
    scheme(iprot).read(iprot, this);
  }

  public void write(org.apache.thrift.protocol.TProtocol oprot) throws org.apache.thrift.TException {
    scheme(oprot).write(oprot, this);
  }

  @Override
  public java.lang.String toString() {
    java.lang.StringBuilder sb = new java.lang.StringBuilder("JMXAttribute(");
    boolean first = true;

    sb.append("attribute:");
    if (this.attribute == null) {
      sb.append("null");
    } else {
      sb.append(this.attribute);
    }
    first = false;
    if (!first) sb.append(", ");
    sb.append("value:");
    if (this.value == null) {
      sb.append("null");
    } else {
      sb.append(this.value);
    }
    first = false;
    sb.append(")");
    return sb.toString();
  }

  public void validate() throws org.apache.thrift.TException {
    // check for required fields
    // check for sub-struct validity
    if (value != null) {
      value.validate();
    }
  }

  private void writeObject(java.io.ObjectOutputStream out) throws java.io.IOException {
    try {
      write(new org.apache.thrift.protocol.TCompactProtocol(new org.apache.thrift.transport.TIOStreamTransport(out)));
    } catch (org.apache.thrift.TException te) {
      throw new java.io.IOException(te);
    }
  }

  private void readObject(java.io.ObjectInputStream in) throws java.io.IOException, java.lang.ClassNotFoundException {
    try {
      read(new org.apache.thrift.protocol.TCompactProtocol(new org.apache.thrift.transport.TIOStreamTransport(in)));
    } catch (org.apache.thrift.TException te) {
      throw new java.io.IOException(te);
    }
  }

  private static class JMXAttributeStandardSchemeFactory implements org.apache.thrift.scheme.SchemeFactory {
    public JMXAttributeStandardScheme getScheme() {
      return new JMXAttributeStandardScheme();
    }
  }

  private static class JMXAttributeStandardScheme extends org.apache.thrift.scheme.StandardScheme<JMXAttribute> {

    public void read(org.apache.thrift.protocol.TProtocol iprot, JMXAttribute struct) throws org.apache.thrift.TException {
      org.apache.thrift.protocol.TField schemeField;
      iprot.readStructBegin();
      while (true)
      {
        schemeField = iprot.readFieldBegin();
        if (schemeField.type == org.apache.thrift.protocol.TType.STOP) { 
          break;
        }
        switch (schemeField.id) {
          case 1: // ATTRIBUTE
            if (schemeField.type == org.apache.thrift.protocol.TType.STRING) {
              struct.attribute = iprot.readString();
              struct.setAttributeIsSet(true);
            } else { 
              org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
            }
            break;
          case 2: // VALUE
            if (schemeField.type == org.apache.thrift.protocol.TType.STRUCT) {
              struct.value = new JMXAttributeValue();
              struct.value.read(iprot);
              struct.setValueIsSet(true);
            } else { 
              org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
            }
            break;
          default:
            org.apache.thrift.protocol.TProtocolUtil.skip(iprot, schemeField.type);
        }
        iprot.readFieldEnd();
      }
      iprot.readStructEnd();

      // check for required fields of primitive type, which can't be checked in the validate method
      struct.validate();
    }

    public void write(org.apache.thrift.protocol.TProtocol oprot, JMXAttribute struct) throws org.apache.thrift.TException {
      struct.validate();

      oprot.writeStructBegin(STRUCT_DESC);
      if (struct.attribute != null) {
        oprot.writeFieldBegin(ATTRIBUTE_FIELD_DESC);
        oprot.writeString(struct.attribute);
        oprot.writeFieldEnd();
      }
      if (struct.value != null) {
        oprot.writeFieldBegin(VALUE_FIELD_DESC);
        struct.value.write(oprot);
        oprot.writeFieldEnd();
      }
      oprot.writeFieldStop();
      oprot.writeStructEnd();
    }

  }

  private static class JMXAttributeTupleSchemeFactory implements org.apache.thrift.scheme.SchemeFactory {
    public JMXAttributeTupleScheme getScheme() {
      return new JMXAttributeTupleScheme();
    }
  }

  private static class JMXAttributeTupleScheme extends org.apache.thrift.scheme.TupleScheme<JMXAttribute> {

    @Override
    public void write(org.apache.thrift.protocol.TProtocol prot, JMXAttribute struct) throws org.apache.thrift.TException {
      org.apache.thrift.protocol.TTupleProtocol oprot = (org.apache.thrift.protocol.TTupleProtocol) prot;
      java.util.BitSet optionals = new java.util.BitSet();
      if (struct.isSetAttribute()) {
        optionals.set(0);
      }
      if (struct.isSetValue()) {
        optionals.set(1);
      }
      oprot.writeBitSet(optionals, 2);
      if (struct.isSetAttribute()) {
        oprot.writeString(struct.attribute);
      }
      if (struct.isSetValue()) {
        struct.value.write(oprot);
      }
    }

    @Override
    public void read(org.apache.thrift.protocol.TProtocol prot, JMXAttribute struct) throws org.apache.thrift.TException {
      org.apache.thrift.protocol.TTupleProtocol iprot = (org.apache.thrift.protocol.TTupleProtocol) prot;
      java.util.BitSet incoming = iprot.readBitSet(2);
      if (incoming.get(0)) {
        struct.attribute = iprot.readString();
        struct.setAttributeIsSet(true);
      }
      if (incoming.get(1)) {
        struct.value = new JMXAttributeValue();
        struct.value.read(iprot);
        struct.setValueIsSet(true);
      }
    }
  }

  private static <S extends org.apache.thrift.scheme.IScheme> S scheme(org.apache.thrift.protocol.TProtocol proto) {
    return (org.apache.thrift.scheme.StandardScheme.class.equals(proto.getScheme()) ? STANDARD_SCHEME_FACTORY : TUPLE_SCHEME_FACTORY).getScheme();
  }
}

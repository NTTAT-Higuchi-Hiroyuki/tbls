package schema

// ObjectType はコメント処理対象のデータベースオブジェクトタイプを表す
type ObjectType string

const (
	// ObjectTypeTable テーブル
	ObjectTypeTable ObjectType = "table"
	// ObjectTypeColumn カラム
	ObjectTypeColumn ObjectType = "column"
	// ObjectTypeIndex インデックス
	ObjectTypeIndex ObjectType = "index"
	// ObjectTypeView ビュー
	ObjectTypeView ObjectType = "view"
	// ObjectTypeConstraint 制約
	ObjectTypeConstraint ObjectType = "constraint"
	// ObjectTypeTrigger トリガー
	ObjectTypeTrigger ObjectType = "trigger"
	// ObjectTypeFunction 関数
	ObjectTypeFunction ObjectType = "function"
	// ObjectTypeEnum 列挙型
	ObjectTypeEnum ObjectType = "enum"
)

// CommentData は解析されたコメントデータを統一的に表現する構造体
type CommentData struct {
	// LogicalName 論理名（表示名）
	LogicalName string `json:"name,omitempty" yaml:"name,omitempty"`
	// Description 説明文
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Tags タグ情報（分類、属性等）
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Metadata 追加メタデータ（キー・バリュー形式）
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	// Priority 優先度（数値が小さいほど高優先度）
	Priority int `json:"priority,omitempty" yaml:"priority,omitempty"`
	// Deprecated 非推奨フラグ
	Deprecated bool `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	// Source 元のコメント文字列（デバッグ用）
	Source string `json:"-" yaml:"-"`
}

// IsEmpty コメントデータが空かどうかを判定
func (cd *CommentData) IsEmpty() bool {
	return cd == nil ||
		(cd.LogicalName == "" &&
			cd.Description == "" &&
			len(cd.Tags) == 0 &&
			len(cd.Metadata) == 0)
}

// HasLogicalName 論理名が設定されているかを判定
func (cd *CommentData) HasLogicalName() bool {
	return cd != nil && cd.LogicalName != ""
}

// HasDescription 説明が設定されているかを判定
func (cd *CommentData) HasDescription() bool {
	return cd != nil && cd.Description != ""
}

// GetLogicalNameOrFallback 論理名を取得、未設定の場合はフォールバック値を返す
func (cd *CommentData) GetLogicalNameOrFallback(fallback string) string {
	if cd != nil && cd.LogicalName != "" {
		return cd.LogicalName
	}
	return fallback
}

// GetDescriptionOrFallback 説明を取得、未設定の場合はフォールバック値を返す
func (cd *CommentData) GetDescriptionOrFallback(fallback string) string {
	if cd != nil && cd.Description != "" {
		return cd.Description
	}
	return fallback
}

// HasTag 指定されたタグが含まれているかを判定
func (cd *CommentData) HasTag(tag string) bool {
	if cd == nil || len(cd.Tags) == 0 {
		return false
	}
	for _, t := range cd.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// GetMetadata 指定されたキーのメタデータを取得
func (cd *CommentData) GetMetadata(key string) (string, bool) {
	if cd == nil || cd.Metadata == nil {
		return "", false
	}
	value, exists := cd.Metadata[key]
	return value, exists
}

// SetMetadata メタデータを設定
func (cd *CommentData) SetMetadata(key, value string) {
	if cd == nil {
		return
	}
	if cd.Metadata == nil {
		cd.Metadata = make(map[string]string)
	}
	cd.Metadata[key] = value
}

// Clone CommentDataの複製を作成
func (cd *CommentData) Clone() *CommentData {
	if cd == nil {
		return nil
	}
	
	clone := &CommentData{
		LogicalName: cd.LogicalName,
		Description: cd.Description,
		Priority:    cd.Priority,
		Deprecated:  cd.Deprecated,
		Source:      cd.Source,
	}
	
	// Tagsのコピー
	if len(cd.Tags) > 0 {
		clone.Tags = make([]string, len(cd.Tags))
		copy(clone.Tags, cd.Tags)
	}
	
	// Metadataのコピー
	if len(cd.Metadata) > 0 {
		clone.Metadata = make(map[string]string)
		for k, v := range cd.Metadata {
			clone.Metadata[k] = v
		}
	}
	
	return clone
}

// Merge 他のCommentDataとマージ（thisが優先、空の場合のみotherで補完）
func (cd *CommentData) Merge(other *CommentData) *CommentData {
	if cd == nil {
		if other == nil {
			return nil
		}
		return other.Clone()
	}
	if other == nil {
		return cd.Clone()
	}
	
	merged := cd.Clone()
	
	// 空の場合のみotherから補完
	if merged.LogicalName == "" && other.LogicalName != "" {
		merged.LogicalName = other.LogicalName
	}
	if merged.Description == "" && other.Description != "" {
		merged.Description = other.Description
	}
	
	// Tagsはマージ（重複除去）
	if len(other.Tags) > 0 {
		tagSet := make(map[string]bool)
		for _, tag := range merged.Tags {
			tagSet[tag] = true
		}
		for _, tag := range other.Tags {
			if !tagSet[tag] {
				merged.Tags = append(merged.Tags, tag)
				tagSet[tag] = true
			}
		}
	}
	
	// Metadataはマージ（thisが優先）
	if len(other.Metadata) > 0 {
		if merged.Metadata == nil {
			merged.Metadata = make(map[string]string)
		}
		for k, v := range other.Metadata {
			if _, exists := merged.Metadata[k]; !exists {
				merged.Metadata[k] = v
			}
		}
	}
	
	return merged
}
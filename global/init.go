package global

// Init 全局的初始化入口
func Init() error {
	var err error

	//todo，初始化配置文件

	// 初始化clientset
	if err = initK8sClient(); err != nil {
		return err
	}
	return nil
}

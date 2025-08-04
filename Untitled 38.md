![image-20230718103904835](/Users/haozhipeng/Library/Application Support/typora-user-images/image-20230718103904835.png)



![1361689647739_.pic_hd](/Users/haozhipeng/Library/Containers/com.tencent.xinWeChat/Data/Library/Application Support/com.tencent.xinWeChat/2.0b4.0.9/7f0e2f5d56dbdecc9c7d740c00231941/Message/MessageTemp/9e20f478899dc29eb19741386f9343c8/Image/1361689647739_.pic_hd.jpg)



1.先执行CBIG_DiffProc_preprocess_tractography.sh

sub2 iFOD2 

```
sh CBIG_DiffProc_preprocess_tractography.sh sub_2 iFOD2 /home/lt292794642/Standalone_CBIG2022_DiffProc-main/stable_projects/preprocessing/CBIG2022_DiffProc/MRtrix /home/lt292794642/input_dir /home/lt292794642/output_dir AMICO /home/lt292794642/output_dir/iFOD2/output/b0_mask/sub_2 "5M" "Schaefer2018_400Parcels_17Networks" "400"

```

2.CBIG_DiffProc_preprocess_tractography.sh中会调用unities下的 generate_b0mask.sh生成b0掩码

```shell
bash CBIG_DiffProc_generate_b0mask.sh /home/lt292794642/input_dir/diffusion sub_1 /home/lt292794642/output_dir/iFOD2/output/b0_mask
```

3.接着会调用gen_parc.sh生成图谱的mgh

4.然后会调用2_gen_tractogram.sh 但需要把scratch删除

```shell
bash CBIG_DiffProc_tractography_2_gen_tractogram.sh sub_1 /home/lt292794642/output_dir/iFOD2/output/sub_1 /home/lt292794642/input_dir/diffusion/sub_1 iFOD2 /home/lt292794642/output_dir/iFOD2/output/b0_mask/sub_1/sub_1_bet_b0_mask.nii.gz 5M
```

5.最后调用3_gen_connectome.sh

```shell
bash CBIG_DiffProc_tractography_3_gen_connectome.sh sub_1 iFOD2 /home/lt292794642/input_dir /home/lt292794642/output_dir 5M Schaefer2018_100Parcels_17Networks 100
```



```shell
#!/bin/bash
#####
# Example: 
#	$CBIG_CODE_DIR/stable_projects/preprocessing/CBIG2022_DiffProc/ \
#		MRtrix/CBIG_DiffProc_batch_tractography.sh \
#		--subj_list /path/to/txtfile --dwi_dir /path/to/dwi_images \
#		--output_dir /path/to/output --py_env name_of_AMICO_environment \
#		--mask_output_dir /path/to/stored/b0_brainmask
#           diff_dir=$input_dir/diffusion/$sub
#           sub
#           tract_streamlines="5M"
#           output_dir /home/lt292794642/output_dir
#           algo = iFOD2
#           sub_outdir= /home/lt292794642/output_dir/iFOD2/output/sub_1
#           mask_output_dir morenweikong zhihouchuangjiancheng /home/lt292794642/output_dir/iFOD2/output/b0_mask/sub_1
#           mask /home/lt292794642/output_dir/iFOD2/output/b0_mask/sub_1/sub_1_bet_b0_mask.nii.gz
#           parcellation_arr="Schaefer2018_100Parcels_17Networks,Schaefer2018_400Parcels_17Networks"
#	    parcels_num_arr="100,400"
# This function fits a tractogram and generates a structural connectivity matrix for each subject
# bash CBIG_DiffProc_generate_b0mask.sh /home/lt292794642/input_dir/diffusion sub_1 /home/lt292794642/output_dir/iFOD2/output/b0_mask
# bash CBIG_DiffProc_tractography_2_gen_tractogram.sh sub_1 /home/lt292794642/output_dir/iFOD2/output/sub_1 /home/lt292794642/input_dir/diffusion/sub_1 iFOD2 /home/lt292794642/output_dir/iFOD2/output/b0_mask/sub_1/sub_1_bet_b0_mask.nii.gz 5M
# bash CBIG_DiffProc_tractography_3_gen_connectome.sh sub_1 iFOD2 /home/lt292794642/input_dir /home/lt292794642/output_dir 5M Schaefer2018_100Parcels_17Networks 100

#
# Written by Leon Ooi and CBIG under MIT license: https://github.com/ThomasYeoLab/CBIG/blob/master/LICENSE.md
#####

###############
# set up environment
#bash CBIG_DiffProc_generate_b0mask.sh /home/lt292794642/input_dir/diffusion sub_1 /home/lt292794642/output_dir/iFOD2/output/b0_mask

###############

# set up directories and lists
scriptdir=/home/lt292794642/Standalone_CBIG2022_DiffProc-main/stable_projects/preprocessing/CBIG2022_DiffProc/MRtrix

while [[ $# -gt 0 ]]; do
key="$1"

case $key in

	-s|--subj_list)
	subj_list="$2"
    	shift; shift;;

	-i|--input_dir)
	input_dir="$2"
    	shift; shift;;

	-o|--output_dir)
	output_dir="$2"
    	shift; shift;;

	-p|--py_env)
	py_env="$2"
    	shift; shift;;

	-m|--mask_output_dir)
	mask_output_dir="$2"
    	shift; shift;;

	-a|--algo)
	algo="$2"
    	shift; shift;;

	-t|--tract_streamlines)
	tract_streamlines="$2"
    	shift; shift;;

	*)    # unknown option
	echo "Unknown option: $1"
	shift;;
esac
done

# check args
if [ -z "$subj_list" ] || [ -z "$input_dir" ] || [ -z "$output_dir" ] || [ -z "$py_env" ]; then
	echo "Missing compulsory variable!"
	exit
fi

if [ -z "$mask_output_dir" ]; then
	mask_output_dir="NIL"
fi

if [ -z "$algo" ]; then
	algo='iFOD2'
fi

if [ -z "$tract_streamlines" ]; then
	tract_streamlines="5M"
fi

# print selected options
echo "[subj_list] = $subj_list"
echo "[input dir] = $input_dir"
echo "[output dir] = $output_dir"
echo "[py env] = $py_env"
echo "[mask_output_dir] = $mask_output_dir"
echo "[algo] = $algo"
echo "[tract_streamlines] = $tract_streamlines"

logdir=$output_dir/$algo/logs
if [ ! -d $logdir ]; then mkdir -p $logdir; fi

###############
# 1. generate required files
###############
# set up list of parcellations
parcellation_arr="Schaefer2018_100Parcels_17Networks,Schaefer2018_400Parcels_17Networks"
parcels_num_arr="100,400"

cat $subj_list | while read subject; do
	# if subject name needs to be modified (e.g. remove prefix), it should be done here
	# for example, delete underscore : sub=$( echo $subject | tr -d '_')
	sub=$( echo $subject )
	echo "[PROGRESS]: Generating parcellations and tractogram for: $sub"
	$scriptdir/CBIG_DiffProc_preprocess_tractography.sh $sub $algo $scriptdir $input_dir $output_dir $py_env \
		$mask_output_dir $tract_streamlines $parcellation_arr $parcels_num_arr
done

###############
# 1.5 wait for job completion
###############
echo "[PROGRESS]: Waiting for all jobs to complete before moving on..."
start=$SECONDS
user_id=$( whoami )
parc_jobs=$( ssh headnode "qselect -u $user_id -N gen_parcellation | wc -l" )
trac_jobs=$( ssh headnode "qselect -u $user_id -N gen_tractogram | wc -l" )
while [[ $parc_jobs -gt 0 ]] || [[ $trac_jobs -gt 0 ]]; do
	duration_min=$((( $SECONDS - $start ) / 60))
	echo "	${duration_min}m: $parc_jobs parcellation jobs & $trac_jobs tractogram jobs still running"
	sleep 3m
	parc_jobs=$( ssh headnode "qselect -u $user_id -N gen_parcellation | wc -l" )
	trac_jobs=$( ssh headnode "qselect -u $user_id -N gen_tractogram | wc -l" )
done

###############
# 2. generate connectome
###############
cat $subj_list | while read subject; do
	# if subject name needs to be modified (e.g. remove prefix), it should be done here
	# for example, delete underscore : sub=$( echo $subject | tr -d '_')
	sub=$( echo $subject )
	echo "[PROGRESS]: Generating connectome for subj ID: $sub"
	cmd="$scriptdir/CBIG_DiffProc_tractography_3_gen_connectome.sh $sub $algo $input_dir $output_dir \
		$tract_streamlines $parcellation_arr $parcels_num_arr > $logdir/${sub}_SC_log.txt 2>&1"
	ssh headnode "source activate $py_env; \
	$CBIG_CODE_DIR/setup/CBIG_pbsubmit -cmd '$cmd' -walltime 02:30:00 \
	-name 'gen_connectome' -mem 6GB -joberr '$logdir' -jobout '$logdir'" < /dev/null
done

###############
# 2.5 wait for job completion
###############
echo "[PROGRESS]: Waiting for all jobs to complete before moving on..."
start=$SECONDS
user_id=$( whoami )
conn_jobs=$( ssh headnode "qselect -u $user_id -N gen_connectome | wc -l" )
while [[ $conn_jobs -gt 0 ]]; do
	duration_min=$((( $SECONDS - $start ) / 60))
	echo "	${duration_min}m: $conn_jobs connectome jobs still running"
	sleep 3m
	conn_jobs=$( ssh headnode "qselect -u $user_id -N gen_connectome | wc -l" )
done

###############
# 3. check logs and replace missing parcels with NaN
###############
cat $subj_list | while read subject; do
	# if subject name needs to be modified (e.g. remove prefix), it should be done here
	# for example, delete underscore : sub=$( echo $subject | tr -d '_')
	sub=$( echo $subject )
	echo "[PROGRESS]: Checking connectomes and filling missing nodes with NaN for subj ID: $sub"
	matlab -nodesktop -nosplash -nodisplay -r " try addpath('$scriptdir'); \
		CBIG_DiffProc_mrtrix_fill_na('$logdir/${sub}_SC_log.txt'); catch ME; \
		display(ME.message); end; exit; " >> $logdir/${sub}_SC_log.txt 2>&1
done

```



```shell
# generate connectome
###############
for parcellation in ${parcellation_arr[@]}; do
    echo "[PROGRESS]: Current connectome = [ $sub_outdir/connectomes/connectome_${parcellation}_SIFT2.csv ]"
    tck2connectome -tck_weights_in $sub_outdir/dwi_wm_weights.csv -assignment_radial_search 2 \
        -stat_edge sum -symmetric $sub_outdir/${tract_streamlines}.tck \
        $sub_outdir/${parcellation}/${parcellation}.mif \
        $sub_outdir/connectomes/connectome_${parcellation}_SIFT2.csv 
    echo "[PROGRESS]: Current connectome = [ $sub_outdir/connectomes/connectome_${parcellation}_unfiltered.csv ]"
    tck2connectome -assignment_radial_search 2 -stat_edge sum -symmetric \
        $sub_outdir/${tract_streamlines}.tck $sub_outdir/${parcellation}/${parcellation}.mif \
        $sub_outdir/connectomes/connectome_${parcellation}_unfiltered.csv -force
    echo "[PROGRESS]: Current connectome = [ $sub_outdir/connectomes/connectome_${parcellation}_length.csv ]"
    tck2connectome -assignment_radial_search 2 -stat_edge mean -symmetric -scale_length \
        $sub_outdir/${tract_streamlines}.tck $sub_outdir/${parcellation}/${parcellation}.mif \
        $sub_outdir/connectomes/connectome_${parcellation}_length.csv 
done

###############
# extract DTI indices
###############
dwi2tensor $sub_outdir/DWI.mif $sub_outdir/DTI_tensor.mif
tensor2metric $sub_outdir/DTI_tensor.mif -fa $sub_outdir/FA.mif -adc $sub_outdir/MD.mif \
    -ad $sub_outdir/AD.mif -rd $sub_outdir/RD.mif
dti_arr=("FA" "MD" "AD" "RD")
for metric in ${dti_arr[@]}; do
    tcksample -stat_tck mean $sub_outdir/${tract_streamlines}.tck $sub_outdir/$metric.mif $sub_outdir/$metric.csv
    if [ ! -d $sub_outdir/connectomes/$metric ]; then mkdir -p $sub_outdir/connectomes/$metric; fi
    for parcellation in ${parcellation_arr[@]}; do
        echo "[PROGRESS]: Current connectome = [ $sub_outdir/connectomes/$metric/connectome_${parcellation}_$metric.csv ]"
        tck2connectome -scale_file $sub_outdir/$metric.csv -stat_edge mean -symmetric \
        $sub_outdir/${tract_streamlines}.tck $sub_outdir/${parcellation}/${parcellation}.mif \
        $sub_outdir/connectomes/$metric/connectome_${parcellation}_$metric.csv 
    done
done
```

首先，使用循环语句遍历parcellation_arr数组中的每个元素（假设该数组包含要处理的分区数据）。

在生成连接组的部分，使用"tck2connectome"命令将跟踪流线（tract streamlines）文件（${tract_streamlines}.tck）与分区文件（${parcellation}/${parcellation}.mif）进行连接，并生成连接组文件（connectome_${parcellation}_SIFT2.csv）。这个命令使用了一些选项，如-tck_weights_in用于指定权重文件，-assignment_radial_search用于设置分配半径搜索距离，-stat_edge用于设置边缘统计方式（这里使用sum），-symmetric用于生成对称的连接组。

接下来，另一个类似的命令生成**未经过滤的连接组文件**（connectome_${parcellation}_unfiltered.csv），仅改变了输出文件的名称。

最后，又一个类似的命令生成带有**长度标度的连接组文件**（connectome_${parcellation}_length.csv），使用了-scale_length选项和-mean统计方式。

在提取DTI指标的部分，首先使用"dwi2tensor"命令将DWI图像（DWI.mif）转换为张量图像（DTI_tensor.mif）。

然后，使用"tensor2metric"命令从张量图像中提取各种DTI指标，包括FA、MD、AD和RD，分别生成相应的图像文件（FA.mif、MD.mif、AD.mif、RD.mif）。

接着，使用"tcksample"命令从跟踪流线文件中提取每个DTI指标的平均值，并将结果保存到相应的.csv文件中。

接下来，对于每个DTI指标，使用嵌套的循环遍历parcellation_arr数组中的每个元素，使用"tck2connectome"命令将DTI指标文件（$sub_outdir/$metric.csv）与分区文件进行连接，并生成相应的连接组文件（connectome_${parcellation}_$metric.csv）。

这段代码的目的是生成连接组文件，以及提取和计算与分区相关的DTI指标。
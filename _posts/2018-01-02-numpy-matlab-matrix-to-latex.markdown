---
layout: post
title:  Convert Matlab - NumPy matrices to LaTeX tables
date:   2018-01-02
author: Paschalis Ts
tags:   [tutorial, wip]
mathjax: false
description: "I'm still getting used to writing *2018*"  


---

### What's this about?
LaTeX is...divisive. It's simultaneously beautiful, frustrating, magnificent and at times, *truly horrendous*. 

This semester, I have been doing some Numerical Analysis for my MSc program, using both Matlab and Numpy.   
Well, here's how I convert matrices from these two environment, into LaTeX tabular/matrix/table to integrate into my reports.   
Other solutions also exist, but who would prefer using external dependencies and searching the internet than hacking together less than 25 lines? Certainly not someone that *chose* to use LaTeX in the first place!  ^^ 

Here are the two gists, one for the [Matlab version](https://gist.github.com/tpaschalis/841dd4a57434ea34506c4408b13547c5) and one for the [Python Version](https://gist.github.com/tpaschalis/7a2943c2248b78b2558c457428086082)


### Where's the sauce?
Feel free to check/fork/use the code, or use a saner method of writing your reports. Well, in any case, here's how it works.

```matlab
function mat2lat(A)

filename = 'tabular.tex';
fileID = fopen(filename, 'w+');
[rows, cols] = size(A);

% Change alignment of your output with the following character
tabalign = repmat('c', 1, cols);

fprintf(fileID,'\\begin{table}[h] \n');
fprintf(fileID,'\\centering \n');
{% raw  %}
fprintf(fileID,'\\begin{tabular}{%s} \n', tabalign);
{% endraw %}
% Change formatting of your output with the following line
tabformat = repmat('%.2f & \t', 1, cols);	% Define output format
tabformat = tabformat(1:end-4);			% Remove last '&' char
tabformat = [tabformat ' \\\\ \n'];		% Add newlines
fprintf(fileID, tabformat, A);			% Spit out input matrix

fprintf(fileID,'\\end{tabular}\n');
fprintf(fileID,'\\end{table}\n');
fclose(fileID);
fprintf('Done :)\n');

end
```


```python
import numpy as np

def np2lat(A):
	filename = 'table.tex'
	f = open(filename, 'a')
	cols = A.shape[1]

	# Change alignment and format of your output
	tabformat = '%.3f'
	tabalign = 'c'*cols

	f.write('\n\\begin{table}[h]\n')
	f.write('\\centering\n')
	{% raw  %}
	f.write('\\begin{tabular}{%s}\n' %tabalign)
	{% endraw %}
	# Use some numpy magic, just addding correct delimiter and newlines
	np.savetxt(f, A, fmt=tabformat, delimiter='\t&\t', newline='\t \\\\ \n')

	f.write('\\end{tabular}\n')
	f.write('\\end{table}\n')
 
M = np.array([[12, 5, 2], [20, 4, 8], [ 2, 4, 3], [ 7, 1, 10]])
np2lat(M)
```


((nil . (
         (eval . (progn
                   (add-to-list 'auto-mode-alist '(".*\\.svelte\\'" . web-mode))
                   (add-to-list 'exec-path (concat (locate-dominating-file default-directory ".dir-locals.el") "node_modules/.bin/"))))
         )))
